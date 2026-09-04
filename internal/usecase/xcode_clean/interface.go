package xcode_clean

import (
	"context"

	"cleaner/internal/domain/entities"
)

// Cleaner — очистка DerivedData Xcode (~/Library/Developer/Xcode/DerivedData).
type Cleaner interface {
	Status(ctx context.Context) (entities.XcodeDerivedDataStatus, error)
	Clean(ctx context.Context) entities.XcodeDerivedDataCleanResult
	// Available сообщает, есть ли на этом устройстве каталог DerivedData.
	Available(ctx context.Context) bool
}