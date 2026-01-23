package ydb

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/ydb-platform/ydb-go-sdk/v3"

	"fmt"
)

func InitNative(ctx context.Context, cfg Config) (*ydb.Driver, error) {
	err := cfg.validate()
	if err != nil {
		return nil, err
	}

	driver, err := ydb.Open(
		ctx,
		cfg.DSN(),
		ydb.WithAccessTokenCredentials(
			cfg.Token,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to open YDB connection: %w", err)
	}

	slog.Info("ydb driver initialized")

	return driver, nil
}

func initDB(ctx context.Context, cfg Config) (*sql.DB, error) {
	driver, err := InitNative(ctx, cfg)
	if err != nil {
		return nil, err
	}

	connector, err := ydb.Connector(driver,
		ydb.WithDefaultQueryMode(ydb.ScriptingQueryMode),
		ydb.WithFakeTx(ydb.ScriptingQueryMode),
		ydb.WithAutoDeclare(),
		ydb.WithNumericArgs(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create YDB connector: %w", err)
	}

	return sql.OpenDB(connector), nil
}
