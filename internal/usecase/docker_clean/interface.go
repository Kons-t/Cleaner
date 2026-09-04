package docker_clean

import (
	"context"

	"cleaner/internal/domain/entities"
)

// Cleaner — очистка неиспользуемых ресурсов Docker (остановленные контейнеры, висячие
// образы, неиспользуемые сети, build cache). Именованные тома не трогает.
type Cleaner interface {
	Status(ctx context.Context) (entities.DockerStatus, error)
	Clean(ctx context.Context) (entities.DockerCleanResult, error)
	// Available сообщает, установлен ли на этом устройстве Docker CLI.
	Available(ctx context.Context) bool
}