package logs_clean

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"

	"cleaner/internal/domain/entities"
	"cleaner/internal/domain/paths"
	"cleaner/internal/pkg/humanize"
)

type usecase struct {
	root string
}

// New создаёт usecase очистки логов приложений. Корень поиска берётся из domain/paths.
func New() Cleaner {
	return &usecase{root: paths.LogsDir}
}

func (u *usecase) Status(ctx context.Context) (entities.LogsStatus, error) {
	files, err := u.findFiles()
	if err != nil {
		return entities.LogsStatus{}, err
	}

	var total int64
	for _, f := range files {
		if info, statErr := os.Stat(f); statErr == nil {
			total += info.Size()
		}
	}

	return entities.LogsStatus{
		Root:      u.root,
		FileCount: len(files),
		SizeHuman: humanize.Bytes(total),
	}, nil
}

func (u *usecase) Clean(ctx context.Context) entities.LogsCleanResult {
	files, err := u.findFiles()
	if err != nil {
		return entities.LogsCleanResult{Errs: []error{err}}
	}

	result := entities.LogsCleanResult{}
	for _, f := range files {
		if rmErr := os.Remove(f); rmErr != nil {
			result.Errs = append(result.Errs, rmErr)
			continue
		}
		result.DeletedCount++
	}

	return result
}

func (u *usecase) Available(ctx context.Context) bool {
	info, err := os.Stat(u.root)

	return err == nil && info.IsDir()
}

// findFiles ищет все файлы (не каталоги) внутри ~/Library/Logs.
func (u *usecase) findFiles() ([]string, error) {
	var files []string

	err := filepath.WalkDir(u.root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			return nil
		}

		files = append(files, path)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}