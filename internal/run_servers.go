package internal

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"time"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"go-ddd-template/internal/application"
	deliveriesinfra "go-ddd-template/internal/infrastructure/deliveries"
	ordersinfra "go-ddd-template/internal/infrastructure/orders"
	productsinfra "go-ddd-template/internal/infrastructure/products"
	usersinfra "go-ddd-template/internal/infrastructure/users"
	productsservice "go-ddd-template/internal/service/products"
	"go-ddd-template/pkg/db/postgres"
	"go-ddd-template/pkg/db/redis"
	pkgydb "go-ddd-template/pkg/db/ydb"
	"go-ddd-template/pkg/logger"
	loggerutils "go-ddd-template/pkg/logger/utils"
	"go-ddd-template/pkg/sqs"
)

//nolint:cyclop
func RunServers(cfg Config, imageStorage productsservice.ImageStorage) error {
	ctx := context.Background()

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

	redisClient, err := redis.NewClient(ctx, cfg.Redis)
	if err != nil {
		return fmt.Errorf("failed to init redis client: %w", err)
	}

	defer func() {
		if err := redisClient.Close(); err != nil {
			slog.Error("Failed to close redis client", loggerutils.ErrAttr(err))
		}
	}()

	driver, err := pkgydb.InitNative(ctx, cfg.YDB)
	if err != nil {
		return fmt.Errorf("failed to init ydb driver: %w", err)
	}

	defer func(ctx context.Context) {
		if err := driver.Close(ctx); err != nil {
			slog.Error("failed to close ydb driver", loggerutils.ErrAttr(err))
		}
	}(ctx)

	sqsSession, err := sqs.NewSession(cfg.SQS)
	if err != nil {
		return fmt.Errorf("failed to init sqs session: %w", err)
	}

	orderQueueWriters, productQueueWriters, queueWritersCloser, err := InitQueueWriters(driver, sqsSession)
	if err != nil {
		return fmt.Errorf("failed to init topic writers: %w", err)
	}

	defer queueWritersCloser(context.Background())

	g, ctx := errgroup.WithContext(ctx)

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	if err := startServers(ctx, g, cluster, orderQueueWriters, productQueueWriters, redisClient, imageStorage, cfg); err != nil {
		return fmt.Errorf("failed to start servers: %w", err)
	}

	if err := g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return fmt.Errorf("server exited with error: %w", err)
	}

	return nil
}

func startServers(
	ctx context.Context,
	g *errgroup.Group,
	cluster *postgres.Cluster,
	queueWriters ordersinfra.QueueWriters,
	productQueueWriters productsinfra.QueueWriters,
	redisClient *redis.Client,
	imageStorage productsservice.ImageStorage,
	cfg Config,
) error {
	usersRepo := usersinfra.NewDBRepo(cluster, redisClient)
	productsRepo := productsinfra.NewDBRepo(cluster, productQueueWriters)
	ordersRepo := ordersinfra.NewDBRepo(cluster, queueWriters)
	deliveriesRepo := deliveriesinfra.NewDBRepo(cluster)

	httpAddr := "0.0.0.0:" + cfg.Server.HTTPPort
	grpcAddr := "0.0.0.0:" + cfg.Server.GRPCPort

	grpcServer, err := startGRPCServer(
		g,
		cfg,
		grpcAddr,
		usersRepo,
		ordersRepo,
		productsRepo,
		deliveriesRepo,
		imageStorage,
	)
	if err != nil {
		return fmt.Errorf("failed to start grpc server: %w", err)
	}

	conn, err := grpc.NewClient(
		grpcAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithConnectParams(grpc.ConnectParams{
			MinConnectTimeout: time.Second,
		}),
	)
	if err != nil {
		return fmt.Errorf("failed to create grpc client: %w", err)
	}

	httpServer, err := startHTTPServer(g, cfg, httpAddr, conn)
	if err != nil {
		return fmt.Errorf("failed to start http server: %w", err)
	}

	// Graceful shutdown logic of http and grpc server
	g.Go(func() error {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.Server.InterruptTimeout)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("failed to shutdown http server: %w", err)
		}

		if err := conn.Close(); err != nil {
			return fmt.Errorf("failed to close grpc connection: %w", err)
		}

		grpcServer.GracefulStop()

		return nil
	})

	if cfg.Server.PprofPort != "" {
		startPprofServer(ctx, g, cfg)
	}

	return nil
}

func startGRPCServer(
	g *errgroup.Group,
	cfg Config,
	addr string,
	usersRepo *usersinfra.DBRepo,
	ordersRepo *ordersinfra.DBRepo,
	productsRepo *productsinfra.PostgresRepo,
	deliveryRepo *deliveriesinfra.DBRepo,
	imageStorage productsservice.ImageStorage,
) (*grpc.Server, error) {
	grpcServer, err := application.SetupGRPCServer(
		cfg.Auth,
		cfg.Traces,
		usersRepo,
		ordersRepo,
		productsRepo,
		deliveryRepo,
		imageStorage,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to setup grpc server: %w", err)
	}

	g.Go(func() error {
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			return fmt.Errorf("failed to listen tcp address [%s]: %w", addr, err)
		}

		slog.Info("Starting grpc server at " + addr)

		if err := grpcServer.Serve(lis); err != nil {
			return fmt.Errorf("failed to serve grpc server: %w", err)
		}

		slog.Info("Grpc server shut down gracefully")

		return nil
	})

	return grpcServer, nil
}

func startHTTPServer(
	g *errgroup.Group,
	cfg Config,
	addr string,
	conn *grpc.ClientConn,
) (*http.Server, error) {
	httpServer, err := application.SetupHTTPServer(conn, cfg.Auth, cfg.Traces)
	if err != nil {
		return nil, fmt.Errorf("failed to setup http server: %w", err)
	}

	g.Go(func() error {
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			return fmt.Errorf("failed to listen tcp address [%s]: %w", addr, err)
		}

		slog.Info("Starting http server at " + addr)

		if err := httpServer.Serve(lis); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("failed to serve http server: %w", err)
		}

		slog.Info("Http server shut down gracefully")

		return nil
	})

	return httpServer, nil
}

func startPprofServer(ctx context.Context, g *errgroup.Group, cfg Config) {
	pprofAddress := "0.0.0.0:" + cfg.Server.PprofPort
	pprofServer := &http.Server{Addr: pprofAddress, Handler: http.DefaultServeMux}

	g.Go(func() error {
		slog.Info("Starting pprof server at " + pprofAddress)

		if err := pprofServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}

		slog.Info("Pprof server shut down gracefully")

		return nil
	})
	g.Go(func() error {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.Server.InterruptTimeout)
		defer cancel()

		err := pprofServer.Shutdown(shutdownCtx)
		if err != nil {
			return err
		}

		return nil
	})
}
