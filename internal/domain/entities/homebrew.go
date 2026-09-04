package entities

// HomebrewStatus — превью того, что удалит `brew cleanup` (вывод --dry-run).
type HomebrewStatus struct {
	Preview string
}

// HomebrewCleanResult — вывод реального запуска `brew cleanup`.
type HomebrewCleanResult struct {
	Output string
}