package main

import (
	"context"
	"fmt"
	"os"

	"cleaner/internal/application"
	"cleaner/internal/usecase/davinci_clean"
	"cleaner/internal/usecase/docker_clean"
	"cleaner/internal/usecase/ds_store_clean"
	"cleaner/internal/usecase/go_clean"
	"cleaner/internal/usecase/homebrew_clean"
	"cleaner/internal/usecase/idea_http_clean"
	"cleaner/internal/usecase/logs_clean"
	"cleaner/internal/usecase/xcode_clean"
)

func main() {
	ctx := context.Background()

	app := application.New(
		go_clean.New(),
		idea_http_clean.New(),
		homebrew_clean.New(),
		docker_clean.New(),
		logs_clean.New(),
		ds_store_clean.New(),
		davinci_clean.New(),
		xcode_clean.New(),
	)

	if err := app.Run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
