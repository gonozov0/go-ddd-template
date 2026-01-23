package ydb

import (
	"context"
	"embed"

	"fmt"

	internal "go-ddd-template/pkg/db/internal"
)

// Migrate applies migrations to YDB database.
// targetVersion is the version to migrate to:
// - nil means migrate to the latest version
// - 0 means rollback all migrations
// - 20240919100509 is an example of a specific version to migrate to.
//
//nolint:cyclop
func Migrate(ctx context.Context, cfg Config, fs embed.FS, targetVersion *uint) error {
	db, err := initDB(ctx, cfg)
	if err != nil {
		return fmt.Errorf("failed to init DB: %w", err)
	}
	defer db.Close()

	if err = internal.Migrate(ctx, fs, targetVersion, db, "ydb"); err != nil {
		return fmt.Errorf("failed to migrate: %w", err)
	}

	return nil
}
