package main

import (
	"context"
	"log/slog"
	"os"

	"fmt"

	"go-ddd-template/internal"
	imagestorage "go-ddd-template/pkg/image_storage"
	loggerutils "go-ddd-template/pkg/logger/utils"
)

func main() {
	cfg, err := internal.LoadConfig()
	if err != nil {
		slog.Error("Could not load config", loggerutils.ErrAttr(fmt.Errorf("failed to load config: %w", err)))
		os.Exit(1)
	}

	imageStorage, err := imagestorage.NewClient(context.Background(), cfg.ImageStorage)
	if err != nil {
		slog.Error(
			"Failed to create image storage client",
			loggerutils.ErrAttr(fmt.Errorf("failed to create image storage client: %w", err)),
		)
		os.Exit(1)
	}

	if err := internal.RunServers(cfg, imageStorage); err != nil {
		slog.Error("Failed to run server", loggerutils.ErrAttr(fmt.Errorf("failed to run server: %w", err)))
		os.Exit(1)
	}
}
