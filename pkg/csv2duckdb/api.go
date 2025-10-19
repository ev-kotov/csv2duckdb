package csv2duckdb

import (
	"context"
	"database/sql"

	"github.com/ev-kotov/csv2duckdb/domain"
	"github.com/ev-kotov/csv2duckdb/internal/importer"
)

// Import CSV files into DuckDB and returns database connection for queries.
//
// DuckDB's automatically detects:
// - Headers presence and names
// - Column separators (comma, semicolon, tab, pipe)
// - Data types (INT, VARCHAR, DATE, TIMESTAMP, BOOLEAN, FLOAT)
// - Date and time formats
// - File encoding
//
// Example:
//
//	db, err := csv2duckdb.Import(ctx,
//	    csv2duckdb.WithFile("data.csv", "users"),
//	    csv2duckdb.WithMemoryLimit(8),
//	)
func Import(ctx context.Context, actions ...domain.Action) (*sql.DB, error) {
	return importer.Import(ctx, actions...)
}
