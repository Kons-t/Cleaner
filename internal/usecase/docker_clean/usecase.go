package docker_clean

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"cleaner/internal/domain/entities"
)

type usecase struct{}

// New создаёт usecase очистки Docker.
func New() Cleaner {
	return &usecase{}
}

func (u *usecase) Status(ctx context.Context) (entities.DockerStatus, error) {
	out, err := exec.CommandContext(ctx, "docker", "system", "df").CombinedOutput()
	if err != nil {
		return entities.DockerStatus{}, fmt.Errorf("docker system df: %w: %s", err, strings.TrimSpace(string(out)))
	}

	return entities.DockerStatus{DiskUsage: strings.TrimSpace(string(out))}, nil
}

// Clean выполняет `docker system prune -f`: убирает остановленные контейнеры, висячие
// образы, неиспользуемые сети и build cache. Именованные тома и используемые образы не трогает.
func (u *usecase) Clean(ctx context.Context) (entities.DockerCleanResult, error) {
	out, err := exec.CommandContext(ctx, "docker", "system", "prune", "-f").CombinedOutput()
	if err != nil {
		return entities.DockerCleanResult{}, fmt.Errorf("docker system prune -f: %w: %s", err, strings.TrimSpace(string(out)))
	}

	return entities.DockerCleanResult{Output: strings.TrimSpace(string(out))}, nil
}

func (u *usecase) Available(ctx context.Context) bool {
	_, err := exec.LookPath("docker")

	return err == nil
}