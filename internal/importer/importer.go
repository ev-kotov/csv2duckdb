package importer

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/duckdb/duckdb-go/v2"
	"github.com/ev-kotov/csv2duckdb/domain"
	"github.com/ev-kotov/csv2duckdb/internal/config"
)

// Import CSV files into DuckDB and returns database connection for queries.
func Import(ctx context.Context, actions ...domain.Action) (*sql.DB, error) {
	startTime := time.Now()

	cfg := config.DefaultConfig()
	for _, action := range actions {
		action(cfg)
	}

	if err := config.ValidateConfig(cfg); err != nil {
		return nil, fmt.Errorf("config validation: %w", err)
	}

	db, err := initDB(ctx, cfg.SavePath)
	if err != nil {
		return nil, fmt.Errorf("init data base: %w", err)
	}

	defer func() {
		if err != nil {
			if closeErr := db.Close(); closeErr != nil {
				err = errors.Join(
					fmt.Errorf("main error: %w", err),
					fmt.Errorf("db close error: %w", closeErr),
				)
			}
		}
	}()

	if err := toOptimize(ctx, db, cfg); err != nil {
		return nil, fmt.Errorf("optimize : %w", err)
	}

	if err := importFiles(ctx, db, cfg); err != nil {
		return nil, fmt.Errorf("import files: %w", err)
	}

	if len(cfg.IndexedColumns) != 0 {
		if err = createIndexes(ctx, db, cfg); err != nil {
			return nil, fmt.Errorf("creating indexes: %w", err)
		}
	}

	log.Printf("CSV import completed in %v", time.Since(startTime))

	if cfg.SavePath != "" {
		log.Printf("Database saved to: %s", cfg.SavePath)
	}

	return db, nil
}

func initDB(ctx context.Context, savePath string) (*sql.DB, error) {
	db, err := sql.Open("duckdb", savePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	if pingErr := db.PingContext(ctx); pingErr != nil {
		pingErr = fmt.Errorf("failed to ping db: %w", pingErr)

		if closeErr := db.Close(); closeErr != nil {
			return nil, errors.Join(
				pingErr,
				fmt.Errorf("failed to close db: %w", closeErr),
			)
		}
		return nil, pingErr
	}

	return db, nil
}
