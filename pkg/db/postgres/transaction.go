package postgres

import (
	"context"
	"log/slog"

	"github.com/jmoiron/sqlx"

	"fmt"

	loggerutils "go-ddd-template/pkg/logger/utils"
)

func RunInTx(ctx context.Context, db *sqlx.DB, fn func(tx *sqlx.Tx) error) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	commited := false

	defer func() {
		if commited {
			return
		}

		if err := tx.Rollback(); err != nil {
			slog.Error(
				"failed to rollback transaction",
				loggerutils.ErrAttr(fmt.Errorf("failed to rollback transaction: %w", err)),
			)
		}
	}()

	if err := fn(tx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	commited = true

	return nil
}
