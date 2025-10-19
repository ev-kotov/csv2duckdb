package config

import (
	"fmt"
	"github.com/ev-kotov/csv2duckdb/domain"
)

// DefaultConfig returns configuration with sensible defaults.
func DefaultConfig() *domain.Config {
	return &domain.Config{
		ProgressBar:            true,
		MemoryLimit:            4,     // 4GB default
		SavePath:               "",    // Empty for in-memory database
		PreserveInsertionOrder: false, // Disable for faster imports (default)
	}
}

// ValidateConfig checks if configuration is valid and sets defaults if needed.
func ValidateConfig(config *domain.Config) error {
	if len(config.Tables) == 0 {
		return fmt.Errorf("no files specified for import")
	}

	if config.MemoryLimit <= 0 {
		config.MemoryLimit = 4 // 4GB default
	}

	return nil
}
