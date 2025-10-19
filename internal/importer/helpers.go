package importer

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"

	"github.com/ev-kotov/csv2duckdb/domain"
)

// applyOptimizations configures DuckDB with performance settings.
func applyOptimizations(db *sql.DB, cfg *domain.Config) error {
	// Use MemoryLimit directly as gigabytes
	memoryLimitGB := cfg.MemoryLimit
	if memoryLimitGB < 1 {
		memoryLimitGB = 1 // Minimum 1GB
	}

	queries := []string{
		fmt.Sprintf("SET memory_limit='%dGB'", memoryLimitGB),
		"SET enable_object_cache=true",
		fmt.Sprintf("SET preserve_insertion_order=%t", cfg.PreserveInsertionOrder),
		fmt.Sprintf("SET enable_progress_bar=%t", cfg.ProgressBar),
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("query %s: %w", query, err)
		}
	}
	return nil
}

// importFiles processes all CSV files specified in configuration.
func importFiles(ctx context.Context, db *sql.DB, cfg *domain.Config) error {
	// Group files by table name for handling same table names
	tableGroups := make(map[string][]string)
	for filePath, tableName := range cfg.Tables {
		tableGroups[tableName] = append(tableGroups[tableName], filePath)
	}

	// Process each table group
	for tableName, filePaths := range tableGroups {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		if len(filePaths) == 1 {
			// Single file for table
			log.Printf("Importing %s -> %s", filePaths[0], tableName)
			if err := importSingleFile(ctx, db, filePaths[0], tableName); err != nil {
				return fmt.Errorf("import %s: %w", filePaths[0], err)
			}
		} else {
			// Multiple files for same table - use UNION ALL
			log.Printf("Importing %d files -> %s", len(filePaths), tableName)
			if err := importMultipleFilesToTable(ctx, db, filePaths, tableName); err != nil {
				return fmt.Errorf("import to %s: %w", tableName, err)
			}
		}
	}
	return nil
}

// importSingleFile imports a single CSV file using DuckDB's read_csv_auto.
func importSingleFile(ctx context.Context, db *sql.DB, filePath, tableName string) error {
	query := fmt.Sprintf(`
		CREATE OR REPLACE TABLE %s AS 
		SELECT * FROM read_csv_auto('%s')`,
		tableName, filePath,
	)
	_, err := db.ExecContext(ctx, query)
	return err
}

// importMultipleFilesToTable imports multiple CSV files into single table using UNION ALL.
func importMultipleFilesToTable(ctx context.Context, db *sql.DB, filePaths []string, tableName string) error {
	// Build UNION ALL query for all files
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

// parallelPostImport executes post-import tasks like index creation.
func parallelPostImport(ctx context.Context, db *sql.DB, cfg *domain.Config) {
	var wg sync.WaitGroup

	if len(cfg.IndexedColumns) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			createIndexes(ctx, db, cfg)
		}()
	}

	wg.Wait()
}

// createIndexes creates indexes on specified columns concurrently.
func createIndexes(ctx context.Context, db *sql.DB, cfg *domain.Config) {
	var wg sync.WaitGroup

	for tableName, columns := range cfg.IndexedColumns {
		for _, column := range columns {
			if ctx.Err() != nil {
				return
			}

			wg.Add(1)
			go func(tblName, col string) {
				defer wg.Done()
				createIndex(ctx, db, tblName, col)
			}(tableName, column)
		}
	}

	wg.Wait()
}

// createIndex creates a single index on specified table and column.
func createIndex(ctx context.Context, db *sql.DB, tableName, column string) {
	indexName := fmt.Sprintf("idx_%s_%s", tableName, column)
	query := fmt.Sprintf("CREATE INDEX IF NOT EXISTS %s ON %s (%s)", indexName, tableName, column)

	if _, err := db.ExecContext(ctx, query); err != nil {
		log.Printf("Failed to create index %s: %v", indexName, err)
	}
}
