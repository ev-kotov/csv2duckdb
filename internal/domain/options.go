package domain

// WithProgressBar включает/выключает прогресс-бар
func WithProgressBar(show bool) Action {
	return func(a *arguments) {
		a.progressBar = show
	}
}

// WithAutoDetectScheme включает/выключает автоопределение схемы
func WithAutoDetectScheme(status bool) Action {
	return func(a *arguments) {
		a.autodetect = status
	}
}

// WithSampleStrategy устанавливает стратегию анализа
func WithSampleStrategy(strategy sampleStrategy) Action {
	return func(a *arguments) {
		a.strategy = strategy
	}
}

// WithHeader указывает наличие заголовка
func WithHeader(header bool) Action {
	return func(a *arguments) {
		a.header = header
	}
}

// WithSeparator устанавливает разделитель
func WithSeparator(separator string) Action {
	return func(a *arguments) {
		a.separator = separator
	}
}

// WithDateFormat устанавливает формат дат
func WithDateFormat(format string) Action {
	return func(a *arguments) {
		a.dateFormat = format
	}
}

// WithTimeStampFormat устанавливает формат временных меток
func WithTimeStampFormat(format string) Action {
	return func(a *arguments) {
		a.timeStampFormat = format
	}
}

// WithMemoryLimit устанавливает лимит памяти
func WithMemoryLimit(limit int) Action {
	return func(a *arguments) {
		a.memoryLimit = limit
	}
}

// WithThreadsCount устанавливает количество потоков
func WithThreadsCount(count int) Action {
	return func(a *arguments) {
		a.threadsCount = count
	}
}

// WithIndexedColumns устанавливает колонки для индексации
func WithIndexedColumns(columns ...string) Action {
	return func(a *arguments) {
		a.indexedColumns = columns
	}
}

// WithMergeFiles включает объединение файлов
func WithMergeFiles(merge bool) Action {
	return func(a *arguments) {
		a.mergeFiles = merge
	}
}
