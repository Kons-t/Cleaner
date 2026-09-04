package ds_store_clean

import (
	"context"

	"cleaner/internal/domain/entities"
)

// Cleaner — очистка мусорных .DS_Store файлов Finder.
type Cleaner interface {
	Status(ctx context.Context) (entities.DSStoreStatus, error)
	Clean(ctx context.Context) entities.DSStoreCleanResult
	// Available сообщает, существует ли на этом устройстве корневая директория поиска.
	Available(ctx context.Context) bool
}