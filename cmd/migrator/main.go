package main

import (
	"context"
	"log/slog"
	"os"

	"fmt"

	"go-ddd-template/internal"
	pgmigrations "go-ddd-template/migrations/postgres"
	ydbmigrations "go-ddd-template/migrations/ydb"
	"go-ddd-template/pkg/db/postgres"
	"go-ddd-template/pkg/db/ydb"
	loggerutils "go-ddd-template/pkg/logger/utils"
)

func main() {
	cfg, err := internal.LoadConfig()
	if err != nil {
		slog.Error("Could not load config", loggerutils.ErrAttr(fmt.Errorf("failed to load config: %w", err)))
		os.Exit(1)
	}

	flags, err := parseFlags()
	if err != nil {
		slog.Error("Failed to parse args", loggerutils.ErrAttr(fmt.Errorf("failed to parse args: %w", err)))
		os.Exit(1)
	}

	ctx := context.Background()
	if flags.DB == YDB {
		if err := ydb.Migrate(ctx, cfg.YDB, ydbmigrations.FS, flags.MigrationVersion); err != nil {
			slog.Error("Failed to migrate YDB", loggerutils.ErrAttr(err))
			os.Exit(1)
		}
	} else {
		if err := postgres.Migrate(ctx, cfg.Postgres, pgmigrations.FS, flags.MigrationVersion); err != nil {
			slog.Error("Failed to migrate Postgres", loggerutils.ErrAttr(err))
			os.Exit(1)
		}
	}

	slog.Info("Migration completed")
}
