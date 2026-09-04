package davinci_clean

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

// New создаёт usecase очистки render cache DaVinci Resolve. Корень поиска берётся из domain/paths.
func New() Cleaner {
	return &usecase{root: paths.DavinciCacheDir}
}

func (u *usecase) Status(ctx context.Context) (entities.DavinciCacheStatus, error) {
	files, err := u.findFiles()
	if err != nil {
		return entities.DavinciCacheStatus{}, err
	}

	var total int64
	for _, f := range files {
		if info, statErr := os.Stat(f); statErr == nil {
			total += info.Size()
		}
	}

	return entities.DavinciCacheStatus{
		Root:      u.root,
		FileCount: len(files),
		SizeHuman: humanize.Bytes(total),
	}, nil
}

func (u *usecase) Clean(ctx context.Context) entities.DavinciCacheCleanResult {
	files, err := u.findFiles()
	if err != nil {
		return entities.DavinciCacheCleanResult{Errs: []error{err}}
	}

	result := entities.DavinciCacheCleanResult{}
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
	_, err := os.Stat(paths.DavinciAppPath)

	return err == nil
}

// findFiles ищет все файлы (не каталоги) внутри render cache DaVinci Resolve.
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
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, err
	}

	return files, nil
}