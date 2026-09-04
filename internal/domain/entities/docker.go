package entities

// DockerStatus — текущее использование диска Docker (вывод `docker system df`).
type DockerStatus struct {
	DiskUsage string
}

// DockerCleanResult — вывод реального запуска `docker system prune`.
type DockerCleanResult struct {
	Output string
}