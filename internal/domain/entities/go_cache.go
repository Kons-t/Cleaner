package entities

import "cleaner/internal/domain/enums"

// CacheDir — расположение кеш-директории и её текущий размер.
type CacheDir struct {
	Path      string
	SizeHuman string
}

// CacheStatus — текущее состояние кешей Go перед очисткой.
type CacheStatus struct {
	GoCache    CacheDir
	GoModCache CacheDir
}

// CleanResult — результат очистки одного типа кеша.
type CleanResult struct {
	Type enums.CacheType
	Err  error
}