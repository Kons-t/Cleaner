package entities

// LogsStatus — сводка по файлам логов приложений в ~/Library/Logs.
type LogsStatus struct {
	Root      string
	FileCount int
	SizeHuman string
}

// LogsCleanResult — результат удаления файлов логов.
type LogsCleanResult struct {
	DeletedCount int
	Errs         []error
}