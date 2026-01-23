package db

import (
	"context"
	"database/sql"
	"embed"
	"log/slog"

	"fmt"

	"github.com/pressly/goose/v3"
)

// Migrate applies migrations to database using goose.
// targetVersion is the version to migrate to:
// - nil means migrate to the latest version
// - 0 means rollback all migrations
// - 20240919100509 is an example of a specific version to migrate to.
func Migrate(ctx context.Context, fs embed.FS, targetVersion *uint, db *sql.DB, dialect string) error {
	goose.SetBaseFS(fs)

	if err := goose.SetDialect(dialect); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	return applyMigrations(ctx, db, targetVersion)
}

func applyMigrations(ctx context.Context, db *sql.DB, targetVersion *uint) error {
	if targetVersion == nil {
		if err := goose.Up(db, "."); err != nil {
			return fmt.Errorf("failed to migrate to the latest version: %w", err)
		}

		slog.Info("Migrations applied successfully")

		return nil
	}

	if *targetVersion == 0 {
		if err := goose.DownTo(db, ".", 0); err != nil {
			return fmt.Errorf("failed to rollback all migrations: %w", err)
		}

		slog.Info("All migrations rolled back successfully")

		return nil
	}

	currentVersion, err := goose.GetDBVersionContext(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to get current version: %w", err)
	}

	convertedTargetVersion := int64(*targetVersion)

	if convertedTargetVersion == currentVersion {
		slog.Info("No migrations to apply")

		return nil
	}

	if convertedTargetVersion < currentVersion {
		if err := goose.DownTo(db, ".", convertedTargetVersion); err != nil {
			return fmt.Errorf("failed to rollback to version %d: %w", *targetVersion, err)
		}

		slog.Info("Successfully rolled back to version", "version", targetVersion)

		return nil
	}

	if err := goose.UpTo(db, ".", convertedTargetVersion); err != nil {
		return fmt.Errorf("failed to migrate to version %d: %w", *targetVersion, err)
	}

	slog.Info("Successfully migrated to version", "version", targetVersion)

	return nil
}
