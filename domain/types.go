package domain

// Config holds configuration for CSV import operations.
type Config struct {
	ProgressBar            bool                // Show progress bar during import
	Tables                 map[string]string   // Map of CSV file paths to table names
	MemoryLimit            int                 // DuckDB memory limit in gigabytes
	IndexedColumns         map[string][]string // Map of table names to columns for indexing
	SavePath               string              // Path to save database file (empty for in-memory)
	PreserveInsertionOrder bool                // Maintain original row order (slower but ordered)
}

// Action is a functional option type for configuring import parameters.
type Action func(*Config)
