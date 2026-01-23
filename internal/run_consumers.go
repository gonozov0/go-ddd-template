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
	"time"

	"golang.org/x/sync/errgroup"

	deliveriesapp "go-ddd-template/internal/application/deliveries/consumers"
	productsapp "go-ddd-template/internal/application/products/consumers"
	deliveriesinfra "go-ddd-template/internal/infrastructure/deliveries"
	productsinfra "go-ddd-template/internal/infrastructure/products"
	"go-ddd-template/pkg/db/postgres"
	pkgydb "go-ddd-template/pkg/db/ydb"
	imagestorage "go-ddd-template/pkg/image_storage"
	"go-ddd-template/pkg/logger"
	loggerutils "go-ddd-template/pkg/logger/utils"
	"go-ddd-template/pkg/sqs"
)

type ConsumerService interface {
	Run(ctx context.Context) error
}

//nolint:cyclop
func RunConsumers(cfg Config) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	g, ctx := errgroup.WithContext(ctx)

	if err := logger.Setup(cfg.Logger); err != nil {
		return fmt.Errorf("failed to setup logger: %w", err)
	}

	cluster, err := postgres.NewCluster(ctx, cfg.Postgres)
	if err != nil {
		return fmt.Errorf("failed to init postgres cluster: %w", err)
	}

	defer func() {
		if err := cluster.Close(); err != nil {
			slog.Error("failed to close postgres cluster", loggerutils.ErrAttr(err))
		}
	}()

	driver, err := pkgydb.InitNative(ctx, cfg.YDB)
	if err != nil {
		return fmt.Errorf("failed to init ydb driver: %w", err)
	}

	defer func() {
		if err := driver.Close(context.Background()); err != nil {
			slog.Error("failed to close ydb driver", loggerutils.ErrAttr(err))
		}
	}()

	sqsSession, err := sqs.NewSession(cfg.SQS)
	if err != nil {
		slog.Error("Failed to init sqs session", loggerutils.ErrAttr(err))
	}

	imageStorage, err := imagestorage.NewClient(ctx, cfg.ImageStorage)
	if err != nil {
		slog.Error("Failed to init image storage", loggerutils.ErrAttr(err))
	}

	deliveryQueueReaders, productQueueReaders, queueReadersCloser, err := InitQueueReaders(driver, sqsSession)
	if err != nil {
		return fmt.Errorf("failed to init topic readers: %w", err)
	}

	defer queueReadersCloser(context.Background())

	_, productQueueWriters, queueWritersCloser, err := InitQueueWriters(driver, sqsSession)
	if err != nil {
		return fmt.Errorf("failed to init topic writers: %w", err)
	}

	defer queueWritersCloser(context.Background())

	deliveryRepo := deliveriesinfra.NewDBRepo(cluster)
	deliveryConsumers := deliveriesapp.NewDeliveryConsumers(deliveryRepo)
	deliveryConsumers.Start(ctx, g, deliveryQueueReaders)

	productRepo := productsinfra.NewDBRepo(cluster, productQueueWriters)
	productConsumers := productsapp.NewDeliveryConsumers(productRepo, imageStorage)
	productConsumers.Start(ctx, g, productQueueReaders)

	startConsumersHTTPServer(ctx, g, cfg)

	slog.Info("Starting consumer service")

	if err := g.Wait(); err != nil {
		return fmt.Errorf("consumer service failed: %w", err)
	}

	slog.Info("Consumer service stopped gracefully")

	return nil
}

const (
	readHeaderTimeout = 2 * time.Second
)

func pingHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

func startConsumersHTTPServer(ctx context.Context, g *errgroup.Group, cfg Config) {
	mux := http.NewServeMux()
	mux.Handle("/ping/", pingHandler())

	httpServer := &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	addr := "0.0.0.0:" + cfg.Consumers.HTTPPort

	g.Go(func() error {
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			return fmt.Errorf("failed to listen tcp address [%s]: %w", addr, err)
		}

		slog.Info("Starting consumers http server at " + addr)

		if err := httpServer.Serve(lis); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("failed to serve consumers http server: %w", err)
		}

		slog.Info("Consumers http server shut down gracefully")

		return nil
	})

	// Graceful shutdown logic of http server
	g.Go(func() error {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.Server.InterruptTimeout)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("failed to shutdown consumers http server: %w", err)
		}

		return nil
	})
}
