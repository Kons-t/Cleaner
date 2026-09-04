package homebrew_clean

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"cleaner/internal/domain/entities"
)

type usecase struct{}

// New создаёт usecase очистки кеша Homebrew.
func New() Cleaner {
	return &usecase{}
}

func (u *usecase) Status(ctx context.Context) (entities.HomebrewStatus, error) {
	out, err := exec.CommandContext(ctx, "brew", "cleanup", "--dry-run").CombinedOutput()
	if err != nil {
		return entities.HomebrewStatus{}, fmt.Errorf("brew cleanup --dry-run: %w: %s", err, strings.TrimSpace(string(out)))
	}

	return entities.HomebrewStatus{Preview: strings.TrimSpace(string(out))}, nil
}

func (u *usecase) Clean(ctx context.Context) (entities.HomebrewCleanResult, error) {
	out, err := exec.CommandContext(ctx, "brew", "cleanup").CombinedOutput()
	if err != nil {
		return entities.HomebrewCleanResult{}, fmt.Errorf("brew cleanup: %w: %s", err, strings.TrimSpace(string(out)))
	}

	return entities.HomebrewCleanResult{Output: strings.TrimSpace(string(out))}, nil
}

func (u *usecase) Available(ctx context.Context) bool {
	_, err := exec.LookPath("brew")

	return err == nil
}