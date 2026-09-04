package homebrew_clean

import (
	"context"

	"cleaner/internal/domain/entities"
)

// Cleaner — очистка кеша Homebrew (старые загрузки, старые версии portable-ruby и т.п.).
type Cleaner interface {
	Status(ctx context.Context) (entities.HomebrewStatus, error)
	Clean(ctx context.Context) (entities.HomebrewCleanResult, error)
	// Available сообщает, установлен ли на этом устройстве Homebrew.
	Available(ctx context.Context) bool
}