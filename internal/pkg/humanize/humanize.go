// Package humanize форматирует размеры в байтах в человекочитаемый вид.
package humanize

import "fmt"

// Bytes форматирует размер в байтах как "1.5 MiB", "12.0 KiB" и т.п.
func Bytes(size int64) string {
	const unit = 1024

	if size < unit {
		return fmt.Sprintf("%d B", size)
	}

	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %ciB", float64(size)/float64(div), "KMGTPE"[exp])
}