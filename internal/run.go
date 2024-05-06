package internal

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	_ "net/http/pprof" //nolint:gosec // pprof port is not exposed to the internet
	"os"
	"os/signal"

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

	startServer(ctx, g, cfg)

	if cfg.Server.PprofPort != "" {
		startPprofServer(ctx, g, cfg)
	}

	if err := g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		slog.Error("Server exited with error", "err", err)
		return err
	}

	slog.Info("Server exited gracefully")
	return nil
}

func startServer(ctx context.Context, g *errgroup.Group, cfg Config) {
	address := "0.0.0.0:" + cfg.Server.Port
	server := newServer(cfg)
	g.Go(func() error {
		slog.Info("Starting server at " + address)
		if err := server.Start(address); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})
	g.Go(func() error {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.Server.InterruptTimeout)
		defer cancel()
		return server.Shutdown(shutdownCtx)
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
		return pprofServer.Shutdown(shutdownCtx)
	})
}
