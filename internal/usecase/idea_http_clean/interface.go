package idea_http_clean

import (
	"context"

	"cleaner/internal/domain/entities"
)

// Cleaner — очистка сохранённых .json-ответов HTTP-клиента IDEA/GoLand (.idea/httpRequests).
type Cleaner interface {
	Status(ctx context.Context) (entities.IdeaHTTPCacheStatus, error)
	Clean(ctx context.Context) entities.IdeaHTTPCacheCleanResult
	// Available сообщает, существует ли на этом устройстве корневая директория поиска.
	Available(ctx context.Context) bool
}