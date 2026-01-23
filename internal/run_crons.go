package internal

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/ydb-platform/ydb-go-sdk/v3"
	"golang.org/x/sync/errgroup"

	orderscrons "go-ddd-template/internal/application/orders/crons"
	ordersinfra "go-ddd-template/internal/infrastructure/orders"
	usersinfra "go-ddd-template/internal/infrastructure/users"
	"go-ddd-template/pkg/db/postgres"
	"go-ddd-template/pkg/db/redis"
	pkgydb "go-ddd-template/pkg/db/ydb"
	"go-ddd-template/pkg/logger"
	loggerutils "go-ddd-template/pkg/logger/utils"
	"go-ddd-template/pkg/sqs"
)

func RunCrons(cfg Config) error {
	if err := logger.Setup(cfg.Logger); err != nil {
		return fmt.Errorf("failed to setup logger: %w", err)
	}

	ctx := context.Background()

	resources, err := initCronResources(ctx, cfg)
	if err != nil {
		return err
	}
	defer resources.Close(ctx)

	g, ctx := errgroup.WithContext(ctx)

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	if err := startCrons(ctx, g, resources, cfg); err != nil {
		return fmt.Errorf("failed to start crons: %w", err)
	}

	if err := g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return fmt.Errorf("server exited with error: %w", err)
	}

	slog.Info("All cron jobs shut down gracefully")

	return nil
}

type cronResources struct {
	cluster      *postgres.Cluster
	redisClient  *redis.Client
	driver       *ydb.Driver
	queueWriters ordersinfra.QueueWriters
	topicCloser  func(context.Context)
}

func initCronResources(ctx context.Context, cfg Config) (*cronResources, error) {
	cluster, err := postgres.NewCluster(ctx, cfg.Postgres)
	if err != nil {
		return nil, fmt.Errorf("failed to init postgres cluster: %w", err)
	}

	driver, err := pkgydb.InitNative(ctx, cfg.YDB)
	if err != nil {
		return nil, fmt.Errorf("failed to init ydb driver: %w", err)
	}

	sqsSession, err := sqs.NewSession(cfg.SQS)
	if err != nil {
		slog.Error("Failed to init sqs session", loggerutils.ErrAttr(err))
	}

	orderQueueWriters, _, queueWritersCloser, err := InitQueueWriters(driver, sqsSession)
	if err != nil {
		return nil, fmt.Errorf("failed to init topic writers: %w", err)
	}

	redisClient, err := redis.NewClient(ctx, cfg.Redis)
	if err != nil {
		return nil, fmt.Errorf("failed to init redis client: %w", err)
	}

	return &cronResources{
		cluster:      cluster,
		redisClient:  redisClient,
		driver:       driver,
		queueWriters: orderQueueWriters,
		topicCloser:  queueWritersCloser,
	}, nil
}

func (r *cronResources) Close(ctx context.Context) {
	if r.redisClient != nil {
		if err := r.redisClient.Close(); err != nil {
			slog.Error("Failed to close redis client", loggerutils.ErrAttr(err))
		}
	}

	if r.topicCloser != nil {
		r.topicCloser(ctx)
	}

	if r.driver != nil {
		if err := r.driver.Close(ctx); err != nil {
			slog.Error("failed to close ydb driver", loggerutils.ErrAttr(err))
		}
	}

	if r.cluster != nil {
		if err := r.cluster.Close(); err != nil {
			slog.Error("failed to close postgres cluster", loggerutils.ErrAttr(err))
		}
	}
}

func startCrons(
	ctx context.Context,
	g *errgroup.Group,
	resources *cronResources,
	cfg Config,
) error {
	ordersRepo := ordersinfra.NewDBRepo(resources.cluster, resources.queueWriters)
	usersRepo := usersinfra.NewDBRepo(resources.cluster, resources.redisClient)

	orderCrons := orderscrons.SetupHandlers(
		ordersRepo,
		usersRepo,
		cfg.Crons.HandleCraetedOrdersInterval,
	)

	slog.Info("Starting cron jobs")
	startCronsHTTPServer(ctx, g, cfg)
	orderCrons.Start(ctx, g)

	return nil
}

func startCronsHTTPServer(ctx context.Context, g *errgroup.Group, cfg Config) {
	mux := http.NewServeMux()
	mux.Handle("/ping/", pingHandler())

	httpServer := &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	addr := "0.0.0.0:" + cfg.Crons.HTTPPort

	g.Go(func() error {
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			return fmt.Errorf("failed to listen tcp address [%s]: %w", addr, err)
		}

		slog.Info("Starting crons http server at " + addr)

		if err := httpServer.Serve(lis); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("failed to serve crons http server: %w", err)
		}

		slog.Info("Crons http server shut down gracefully")

		return nil
	})

	// Graceful shutdown logic of http server
	g.Go(func() error {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.Server.InterruptTimeout)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("failed to shutdown crons http server: %w", err)
		}

		return nil
	})
}
