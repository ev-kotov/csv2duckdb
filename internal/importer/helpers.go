package importer

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/ev-kotov/csv2duckdb/domain"
)

func toOptimize(ctx context.Context, db *sql.DB, cfg *domain.Config) error {
	memoryLimit := cfg.MemoryLimit

	if memoryLimit < 1 {
		memoryLimit = 1 // Minimum 1GB
	}

	queries := []string{
		fmt.Sprintf("SET memory_limit='%dGB'", memoryLimit),
		"SET enable_object_cache=true",
		fmt.Sprintf("SET preserve_insertion_order=%t", cfg.PreserveInsertionOrder),
		fmt.Sprintf("SET enable_progress_bar=%t", cfg.ProgressBar),
	}

	for _, query := range queries {
		if _, err := db.ExecContext(ctx, query); err != nil {
			return fmt.Errorf("query %s: %w", query, err)
		}
	}

	return nil
}

func importFiles(ctx context.Context, db *sql.DB, cfg *domain.Config) error {

	tableGroups := make(map[string][]string)
	for filePath, tableName := range cfg.Tables {
		tableGroups[tableName] = append(tableGroups[tableName], filePath)
	}

	for tableName, filePaths := range tableGroups {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		if len(filePaths) == 1 {
			log.Printf("Importing %s -> %s", filePaths[0], tableName)
			if err := importOneFile(ctx, db, filePaths[0], tableName); err != nil {
				return fmt.Errorf("import %s: %w", filePaths[0], err)
			}
		} else {
			log.Printf("Importing %d files -> %s", len(filePaths), tableName)
			if err := importSameFiles(ctx, db, filePaths, tableName); err != nil {
				return fmt.Errorf("import to %s: %w", tableName, err)
			}
		}
	}
	return nil
}

func importOneFile(ctx context.Context, db *sql.DB, filePath, tableName string) error {
	query := fmt.Sprintf(`
		CREATE OR REPLACE TABLE %s AS 
		SELECT * FROM read_csv_auto('%s')`,
		tableName, filePath,
	)
	_, err := db.ExecContext(ctx, query)
	return err
}

func importSameFiles(ctx context.Context, db *sql.DB, filePaths []string, tableName string) error {
	unionQuery := ""
	for i, filePath := range filePaths {
		if i > 0 {
			unionQuery += " UNION ALL "
		}
		unionQuery += fmt.Sprintf("SELECT * FROM read_csv_auto('%s')", filePath)
	}

	query := fmt.Sprintf(`
		CREATE OR REPLACE TABLE %s AS 
		%s`,
		tableName, unionQuery,
	)
	_, err := db.ExecContext(ctx, query)
	return err
}

func createIndexes(ctx context.Context, db *sql.DB, cfg *domain.Config) error {
	var wg sync.WaitGroup

	errCh := make(chan error, len(cfg.IndexedColumns)*10)

	for tableName, columns := range cfg.IndexedColumns {
		for _, column := range columns {

			if ctx.Err() != nil {
				return ctx.Err()
			}

			wg.Add(1)

			go func(tblName, col string) {

				defer wg.Done()

				if err := createIndex(ctx, db, tblName, col); err != nil {
					errCh <- fmt.Errorf("table %s column %s: %w", tblName, col, err)
				}
			}(tableName, column)
		}
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	var errs []error
	for err := range errCh {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return fmt.Errorf("index creation errors: %w", errors.Join())
	}

	return nil
}

func createIndex(ctx context.Context, db *sql.DB, tableName, column string) error {
	indexName := fmt.Sprintf("idx_%s_%s", tableName, column)

	query := fmt.Sprintf("CREATE INDEX IF NOT EXISTS %s ON %s (%s)",
		indexName, tableName, column)

	if _, err := db.ExecContext(ctx, query); err != nil {
		return fmt.Errorf("failed to create index %s: %w", indexName, err)
	}
	return nil
}
