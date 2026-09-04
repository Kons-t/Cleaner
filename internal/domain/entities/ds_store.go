package entities

// DSStoreStatus — сводка по мусорным .DS_Store файлам Finder.
type DSStoreStatus struct {
	Root      string
	FileCount int
	SizeHuman string
}

// DSStoreCleanResult — результат удаления .DS_Store файлов.
type DSStoreCleanResult struct {
	DeletedCount int
	Errs         []error
}