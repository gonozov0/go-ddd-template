package postgres

import (
	"context"
	"embed"

	"fmt"

	internal "go-ddd-template/pkg/db/internal"
)

// Migrate applies migrations to Postgres database with goose.
// targetVersion is the version to migrate to:
// - nil means migrate to the latest version
// - 0 means rollback all migrations
// - 20240919100509 is an example of a specific version to migrate to.
func Migrate(ctx context.Context, cfg Config, fs embed.FS, targetVersion *uint) error {
	cluster, err := NewCluster(ctx, cfg)
	if err != nil {
		return fmt.Errorf("failed to init postgres cluster: %w", err)
	}
	defer cluster.Close()

	db, err := cluster.PrimaryDB(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get primary db: %w", err)
	}

	return internal.Migrate(ctx, fs, targetVersion, db, "postgres")
}
