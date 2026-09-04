package go_clean

import (
	"context"

	"cleaner/internal/domain/entities"
	"cleaner/internal/domain/enums"
)

// Cleaner — очистка кешей Go: показать текущее состояние и очистить выбранные типы.
type Cleaner interface {
	Status(ctx context.Context) (entities.CacheStatus, error)
	Clean(ctx context.Context, types []enums.CacheType) []entities.CleanResult
	// Available сообщает, установлен ли на этом устройстве инструмент `go`.
	Available(ctx context.Context) bool
}