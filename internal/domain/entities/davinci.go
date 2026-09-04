package entities

// DavinciCacheStatus — сводка по файлам render cache DaVinci Resolve.
type DavinciCacheStatus struct {
	Root      string
	FileCount int
	SizeHuman string
}

// DavinciCacheCleanResult — результат удаления файлов render cache.
type DavinciCacheCleanResult struct {
	DeletedCount int
	Errs         []error
}