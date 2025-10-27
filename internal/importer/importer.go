package importer

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/ev-kotov/csv2duckdb/domain"
	"github.com/ev-kotov/csv2duckdb/internal/config"
	_ "github.com/marcboeker/go-duckdb"
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
		return nil, fmt.Errorf("initialization: %w", err)
	}

	if err := toOptimize(ctx, db, cfg); err != nil {
		return nil, fmt.Errorf("optimization: %w", err)
	}

	if err := importFiles(ctx, db, cfg); err != nil {
		return nil, fmt.Errorf("importing: %w", err)
	}

	parallelPostImport(ctx, db, cfg)

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

	if err := db.PingContext(ctx); err != nil {
		err := db.Close()

		if err := db.Close(); err != nil {
			return nil, fmt.Errorf("failed to close db: %w", err)
		}

		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	return db, nil
}
