package logs_clean

import (
	"context"

	"cleaner/internal/domain/entities"
)

// Cleaner — очистка файлов логов приложений в ~/Library/Logs.
type Cleaner interface {
	Status(ctx context.Context) (entities.LogsStatus, error)
	Clean(ctx context.Context) entities.LogsCleanResult
	// Available сообщает, существует ли на этом устройстве директория логов.
	Available(ctx context.Context) bool
}