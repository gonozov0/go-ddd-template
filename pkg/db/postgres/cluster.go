package postgres

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"golang.yandex/hasql/checkers"
	hasql "golang.yandex/hasql/sqlx"

	"fmt"
)

const (
	firstConnectionTimeout = 5 * time.Second
	clusterUpdateInterval  = 5 * time.Second
	clusterUpdateTimeout   = 5 * time.Second
	selectNodeTimeout      = time.Second
)

func initCluster(ctx context.Context, connData Config) (*hasql.Cluster, error) {
	nodes, err := initNodes(connData)
	if err != nil {
		return nil, fmt.Errorf("failed to init nodes: %w", err)
	}

	opts := []hasql.ClusterOption{
		hasql.WithUpdateInterval(clusterUpdateInterval),
		hasql.WithUpdateTimeout(clusterUpdateTimeout),
	}
	cluster, _ := hasql.NewCluster(nodes, checkers.PostgreSQL, opts...)

	ctx2, cancel := context.WithTimeout(ctx, firstConnectionTimeout)
	defer cancel()

	_, err = cluster.WaitForAlive(ctx2)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for cluster availability: %w", err)
	}

	return cluster, nil
}

func initNodes(cfg Config) ([]hasql.Node, error) {
	if len(cfg.Hosts) == 0 {
		return nil, fmt.Errorf("no host found")
	}

	nodes := make([]hasql.Node, 0, len(cfg.Hosts))

	for _, host := range cfg.Hosts {
		dsn, err := pgx.ParseConfig(cfg.GetConnStr(host))
		if err != nil {
			return nil, fmt.Errorf("failed to parse config: %w", err)
		}

		dsn.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

		db := sqlx.NewDb(stdlib.OpenDB(*dsn), "pgx")
		nodes = append(nodes, hasql.NewNode(host, db))
	}

	return nodes, nil
}

type Cluster struct {
	cluster *hasql.Cluster
}

func NewCluster(ctx context.Context, connData Config) (*Cluster, error) {
	cluster, err := initCluster(ctx, connData)
	if err != nil {
		return nil, fmt.Errorf("failed to init cluster: %w", err)
	}

	slog.Info("postgres cluster initialized")

	return &Cluster{cluster: cluster}, nil
}

func (c *Cluster) getPrimaryNode(ctx context.Context) (hasql.Node, error) {
	ctx, cancel := context.WithTimeout(ctx, selectNodeTimeout)
	defer cancel()

	node, err := c.cluster.WaitForPrimary(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for primary: %w", err)
	}

	return node, nil
}

func (c *Cluster) PrimaryDBx(ctx context.Context) (*sqlx.DB, error) {
	node, err := c.getPrimaryNode(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get primary node: %w", err)
	}

	return node.DBx(), nil
}

func (c *Cluster) PrimaryDB(ctx context.Context) (*sql.DB, error) {
	node, err := c.getPrimaryNode(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get primary node: %w", err)
	}

	return node.DB(), nil
}

func (c *Cluster) StandbyPreferredDBx(ctx context.Context) (*sqlx.DB, error) {
	ctx, cancel := context.WithTimeout(ctx, selectNodeTimeout)
	defer cancel()

	node, err := c.cluster.WaitForStandbyPreferred(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for standby preferred: %w", err)
	}

	return node.DBx(), nil
}

func (c *Cluster) Close() error {
	return c.cluster.Close()
}
