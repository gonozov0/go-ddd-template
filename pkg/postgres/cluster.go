package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"golang.yandex/hasql/checkers"
	hasql "golang.yandex/hasql/sqlx"
)

const (
	connectionTimeout     = time.Second
	clusterUpdateInterval = 2 * time.Second
)

func InitCluster(ctx context.Context, connData ConnectionData) (*hasql.Cluster, error) {
	nodes, err := initNodes(ctx, connData)
	if err != nil {
		return nil, fmt.Errorf("failed to init nodes: %w", err)
	}
	opts := []hasql.ClusterOption{hasql.WithUpdateInterval(clusterUpdateInterval)}
	cluster, _ := hasql.NewCluster(nodes, checkers.PostgreSQL, opts...)
	ctx2, cancel := context.WithTimeout(ctx, connectionTimeout)
	defer cancel()
	_, err = cluster.WaitForPrimary(ctx2)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for cluster availability: %w", err)
	}

	return cluster, nil
}

func initNodes(ctx context.Context, connData ConnectionData) ([]hasql.Node, error) {
	nodes := make([]hasql.Node, 0, len(connData.Hosts))
	for _, host := range connData.Hosts {
		dsn, err := pgx.ParseConfig(connData.String(host))
		if err != nil {
			return nil, fmt.Errorf("failed to parse config: %w", err)
		}
		db := sqlx.NewDb(stdlib.OpenDB(*dsn), "pgx")
		if err := pingWithTimeout(ctx, db); err != nil {
			return nil, fmt.Errorf("failed to ping host %s: %w", host, err)
		}
		nodes = append(nodes, hasql.NewNode(host, db))
	}
	return nodes, nil
}

func pingWithTimeout(ctx context.Context, db *sqlx.DB) error {
	ctx2, cancel := context.WithTimeout(ctx, connectionTimeout)
	defer cancel()
	return db.PingContext(ctx2)
}
