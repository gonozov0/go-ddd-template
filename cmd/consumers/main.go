package main

import (
	"log/slog"
	"os"

	"fmt"

	"go-ddd-template/internal"
	loggerutils "go-ddd-template/pkg/logger/utils"
)

func main() {
	if err := run(); err != nil {
		slog.Error("consumer failed", loggerutils.ErrAttr(err))
		os.Exit(1)
	}
}

func run() error {
	cfg, err := internal.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	return internal.RunConsumers(cfg)
}
