package csv2duckdb

import "github.com/ev-kotov/csv2duckdb/domain"

// WithTable specifies a single CSV file and target table name.
// If multiple files use same table name, they will be merged into one table.
func WithTable(filePath, tableName string) domain.Action {
	return func(c *domain.Config) {
		if c.Tables == nil {
			c.Tables = make(map[string]string)
		}
		c.Tables[filePath] = tableName
	}
}

// WithTables specifies multiple CSV files and their target table names.
// Keys are CSV file paths, values are table names.
// Files with same table name will be merged into single table.
func WithTables(files map[string]string) domain.Action {
	return func(c *domain.Config) {
		c.Tables = files
	}
}

// WithProgressBar controls whether to show import progress bar.
// Default: true
func WithProgressBar(show bool) domain.Action {
	return func(c *domain.Config) {
		c.ProgressBar = show
	}
}

// WithMemoryLimit sets the memory limit for DuckDB in gigabytes.
// Example: WithMemoryLimit(4) sets 4GB memory limit
// Default: 4GB
func WithMemoryLimit(gigabytes int) domain.Action {
	return func(c *domain.Config) {
		c.MemoryLimit = gigabytes
	}
}

// WithIndexedColumns specifies columns to create indexes for specific tables.
// Keys are table names, values are slices of column names to index.
// Example: WithIndexedColumns(map[string][]string{"users": {"id", "email"}, "orders": {"user_id"}})
// Default: empty (no indexes)
func WithIndexedColumns(indexes map[string][]string) domain.Action {
	return func(c *domain.Config) {
		c.IndexedColumns = indexes
	}
}

// WithSaveAs specifies the file path to save the database.
// If empty (default), uses in-memory database.
// Example: WithSaveAs("my_database.duckdb") saves to file
// Default: "" (in-memory)
func WithSaveAs(path string) domain.Action {
	return func(c *domain.Config) {
		c.SavePath = path
	}
}

// WithInsertionOrder enables or disables preservation of row insertion order.
// When enabled (true), maintains original CSV row order (slower import).
// When disabled (false), allows DuckDB to optimize import speed (faster import).
// Use true only when row order is critical for your application logic.
// Default: false (optimized for performance)
func WithInsertionOrder(preserve bool) domain.Action {
	return func(c *domain.Config) {
		c.PreserveInsertionOrder = preserve
	}
}
