package go_clean

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"cleaner/internal/domain/entities"
	"cleaner/internal/domain/enums"
	"cleaner/internal/pkg/humanize"
)

type usecase struct{}

// New создаёт usecase очистки кешей Go.
func New() Cleaner {
	return &usecase{}
}

func (u *usecase) Status(ctx context.Context) (entities.CacheStatus, error) {
	goCachePath, err := goEnv(ctx, "GOCACHE")
	if err != nil {
		return entities.CacheStatus{}, err
	}

	goModCachePath, err := goEnv(ctx, "GOMODCACHE")
	if err != nil {
		return entities.CacheStatus{}, err
	}

	var wg sync.WaitGroup
	var goCacheSize, goModCacheSize int64

	wg.Add(2)
	go func() {
		defer wg.Done()
		goCacheSize = dirSize(goCachePath)
	}()
	go func() {
		defer wg.Done()
		goModCacheSize = dirSize(goModCachePath)
	}()
	wg.Wait()

	return entities.CacheStatus{
		GoCache:    entities.CacheDir{Path: goCachePath, SizeHuman: humanize.Bytes(goCacheSize)},
		GoModCache: entities.CacheDir{Path: goModCachePath, SizeHuman: humanize.Bytes(goModCacheSize)},
	}, nil
}

func (u *usecase) Clean(ctx context.Context, types []enums.CacheType) []entities.CleanResult {
	results := make([]entities.CleanResult, 0, len(types))

	for _, t := range types {
		cmd := exec.CommandContext(ctx, "go", "clean", t.GoCleanFlag())
		_, err := cmd.CombinedOutput()

		results = append(results, entities.CleanResult{Type: t, Err: err})
	}

	return results
}

func (u *usecase) Available(ctx context.Context) bool {
	_, err := exec.LookPath("go")

	return err == nil
}

func goEnv(ctx context.Context, key string) (string, error) {
	out, err := exec.CommandContext(ctx, "go", "env", key).Output()
	if err != nil {
		return "", fmt.Errorf("go env %s: %w", key, err)
	}

	return strings.TrimSpace(string(out)), nil
}

// dirSize считает суммарный размер директории в байтах. При ошибке чтения возвращает то, что успел насчитать.
func dirSize(path string) int64 {
	var size int64

	_ = filepath.WalkDir(path, func(_ string, d os.DirEntry, err error) error {
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