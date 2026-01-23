package main

import (
	"log/slog"
	"os"

	"go-ddd-template/internal"
	loggerutils "go-ddd-template/pkg/logger/utils"
)

func main() {
	cfg, err := internal.LoadConfig()
	if err != nil {
		slog.Error("Could not load config", loggerutils.ErrAttr(err))
		os.Exit(1)
	}

	if err := internal.RunCrons(cfg); err != nil {
		slog.Error("Failed to run crons", loggerutils.ErrAttr(err))
		os.Exit(1)
	}
}
