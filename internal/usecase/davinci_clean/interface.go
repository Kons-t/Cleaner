package davinci_clean

import (
	"context"

	"cleaner/internal/domain/entities"
)

// Cleaner — очистка render cache DaVinci Resolve (~/Movies/CacheClip).
type Cleaner interface {
	Status(ctx context.Context) (entities.DavinciCacheStatus, error)
	Clean(ctx context.Context) entities.DavinciCacheCleanResult
	// Available сообщает, установлен ли DaVinci Resolve на этом устройстве.
	Available(ctx context.Context) bool
}