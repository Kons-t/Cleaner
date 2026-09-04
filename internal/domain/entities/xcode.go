package entities

// XcodeDerivedDataStatus — сводка по содержимому DerivedData Xcode.
// EntryCount — число верхнеуровневых папок проектов, а не отдельных файлов.
type XcodeDerivedDataStatus struct {
	Root       string
	EntryCount int
	SizeHuman  string
}

// XcodeDerivedDataCleanResult — результат удаления папок проектов из DerivedData.
type XcodeDerivedDataCleanResult struct {
	DeletedCount int
	Errs         []error
}