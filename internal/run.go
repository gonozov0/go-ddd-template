package internal

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	_ "net/http/pprof" //nolint:gosec // pprof port is not exposed to the internet
	"os"
	"os/signal"
	"strings"

	"github.com/soheilhy/cmux"

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

	err = startServers(ctx, g, cfg)
	if err != nil {
		slog.Error("Could not start servers", "err", err)
		return err
	}

	if cfg.Server.PprofPort != "" {
		startPprofServer(ctx, g, cfg)
	}

	if err := g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		slog.Error("Server exited with error", "err", err)
		return err
	}
	return nil
}

func startServers(ctx context.Context, g *errgroup.Group, cfg Config) error {
	listener, err := net.Listen("tcp", "0.0.0.0:"+cfg.Server.Port)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	m := cmux.New(listener)
	grpcListener := m.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
	httpListener := m.Match(cmux.Any())

	userRepo := usersInfra.NewInMemoryRepo()
	productRepo := productsInfra.NewInMemoryRepo()
	orderRepo := ordersInfra.NewInMemoryRepo()

	httpServer := application.SetupHTTPServer(userRepo, orderRepo, productRepo)
	grpcServer := application.SetupGRPCServer(userRepo, orderRepo, productRepo)

	slog.Info("Starting server http and grpc server 0.0.0.0:" + cfg.Server.Port)
	g.Go(func() error {
		//nolint:gosec,govet // timeouts set up on balancers
		if err := http.Serve(httpListener, httpServer); err != nil {
			if errors.Is(err, cmux.ErrServerClosed) {
				slog.Info("HTTP server closed")
				return nil
			}
			return err
		}
		return nil
	})
	g.Go(func() error {
		//nolint:govet // err is allowed to shadow
		if err := grpcServer.Serve(grpcListener); err != nil {
			if errors.Is(err, cmux.ErrServerClosed) {
				slog.Info("GRPC server closed")
				return nil
			}
			return err
		}
		return nil
	})
	g.Go(func() error {
		err = m.Serve()
		if err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") {
				slog.Info("shutting down a mux server")
				return nil
			}
			slog.Error("Failed to serve", "err", err)
			return err
		}
		return nil
	})
	g.Go(func() error {
		<-ctx.Done()
		m.Close()
		slog.Info("Server shut down gracefully")
		return nil
	})

	return nil
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
