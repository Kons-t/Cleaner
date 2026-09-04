package xcode_clean

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

// New создаёт usecase очистки DerivedData Xcode. Корень поиска берётся из domain/paths.
func New() Cleaner {
	return &usecase{root: paths.XcodeDerivedDataDir}
}

func (u *usecase) Status(ctx context.Context) (entities.XcodeDerivedDataStatus, error) {
	entries, err := u.findEntries()
	if err != nil {
		return entities.XcodeDerivedDataStatus{}, err
	}

	var total int64
	for _, e := range entries {
		total += dirSize(e)
	}

	return entities.XcodeDerivedDataStatus{
		Root:       u.root,
		EntryCount: len(entries),
		SizeHuman:  humanize.Bytes(total),
	}, nil
}

func (u *usecase) Clean(ctx context.Context) entities.XcodeDerivedDataCleanResult {
	entries, err := u.findEntries()
	if err != nil {
		return entities.XcodeDerivedDataCleanResult{Errs: []error{err}}
	}

	result := entities.XcodeDerivedDataCleanResult{}
	for _, e := range entries {
		if rmErr := os.RemoveAll(e); rmErr != nil {
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

// findEntries возвращает пути верхнеуровневых записей (обычно — по одной папке на проект) внутри DerivedData.
func (u *usecase) findEntries() ([]string, error) {
	dirEntries, err := os.ReadDir(u.root)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, err
	}

	entries := make([]string, 0, len(dirEntries))
	for _, e := range dirEntries {
		entries = append(entries, filepath.Join(u.root, e.Name()))
	}

	return entries, nil
}

// dirSize считает суммарный размер файлов внутри path (файл или каталог) в байтах.
// При ошибке чтения возвращает то, что успел насчитать.
func dirSize(path string) int64 {
	var size int64

	_ = filepath.WalkDir(path, func(_ string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return nil
		}

		size += info.Size()

		return nil
	})

	return size
}