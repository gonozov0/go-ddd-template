package internal

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	_ "net/http/pprof" //nolint:gosec // pprof port is not exposed to the internet
	"os"
	"os/signal"
	"strings"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"go-echo-ddd-template/internal/application"
	ordersInfra "go-echo-ddd-template/internal/infrastructure/orders"
	productsInfra "go-echo-ddd-template/internal/infrastructure/products"
	usersInfra "go-echo-ddd-template/internal/infrastructure/users"
	"go-echo-ddd-template/pkg/logger"
	"go-echo-ddd-template/pkg/sentry"

	"golang.org/x/sync/errgroup"
)

func Run() error {
	cfg, err := LoadConfig()
	if err != nil {
		slog.Error("Could not load config", "err", err)
		return err
	}
	if err = sentry.Init(cfg.Sentry.DSN, cfg.Sentry.Environment); err != nil {
		slog.Error("Could not init sentry", "err", err)
		return err
	}
	logger.Setup()

	g, ctx := errgroup.WithContext(context.Background())
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	startServers(ctx, g, cfg)
	if cfg.Server.PprofPort != "" {
		startPprofServer(ctx, g, cfg)
	}

	if err := g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		slog.Error("Server exited with error", "err", err)
		return err
	}
	return nil
}

func startServers(ctx context.Context, g *errgroup.Group, cfg Config) {
	userRepo := usersInfra.NewInMemoryRepo()
	productRepo := productsInfra.NewInMemoryRepo()
	orderRepo := ordersInfra.NewInMemoryRepo()

	httpServer := application.SetupHTTPServer(userRepo, orderRepo, productRepo)
	grpcServer := application.SetupGRPCServer(userRepo, orderRepo, productRepo)

	address := "0.0.0.0:" + cfg.Server.Port
	server := &http.Server{
		Addr: address,
		Handler: h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
				grpcServer.ServeHTTP(w, r)
			} else {
				httpServer.ServeHTTP(w, r)
			}
		}), &http2.Server{}),
		ReadHeaderTimeout: cfg.Server.ReadHeaderTimeout,
	}

	slog.Info("Starting server http and grpc server at " + address)
	g.Go(func() error {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})
	g.Go(func() error {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.Server.InterruptTimeout)
		defer cancel()
		err := httpServer.Shutdown(shutdownCtx)
		if err != nil {
			return err
		}
		slog.Info("Http and grpc server shut down gracefully")
		return nil
	})
}

func startPprofServer(ctx context.Context, g *errgroup.Group, cfg Config) {
	pprofAddress := "0.0.0.0:" + cfg.Server.PprofPort
	//nolint:gosec // pprofServer is not exposed to the internet
	pprofServer := &http.Server{Addr: pprofAddress, Handler: http.DefaultServeMux}
	g.Go(func() error {
		slog.Info("Starting pprof server at " + pprofAddress)
		if err := pprofServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
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
		slog.Info("Pprof server shut down gracefully")
		return nil
	})
}
