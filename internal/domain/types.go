package domain

import "time"

type sampleStrategy int

const (
	SampleAllRows sampleStrategy = -1
	SampleAuto                   = 0
)

// Action implements func for main arguments.
type Action func(*arguments)

type arguments struct {
	progressBar              bool
	filesPathsAndTablesNames map[string]string
	mergeFiles               bool
	autodetect               bool
	strategy                 sampleStrategy
	header                   bool
	separator                string
	dateFormat               string
	timeStampFormat          string
	memoryLimit              int
	threadsCount             int
	indexedColumns           []string
}

// ImportResult результат импорта
type ImportResult struct {
	DB         interface{} // *sql.DB - будет в implementation
	TableName  string
	RowCount   int64
	ImportTime time.Duration
	TableSize  string
}
