package internal

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go-echo-ddd-template/pkg/logger"
	"go-echo-ddd-template/pkg/sentry"

	"golang.org/x/sync/errgroup"
)

func Run() error {
	config, err := LoadConfig()
	if err != nil {
		slog.Error("Could not load config", "err", err)
		return err
	}
	if err = sentry.Init(config.Sentry.DSN, config.Sentry.Environment); err != nil {
		slog.Error("Could not init sentry", "err", err)
		return err
	}
	logger.Setup()

	server := newServer(config)
	g, ctx := errgroup.WithContext(context.Background())
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	g.Go(func() error {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), config.Server.InterruptTimeout)
		defer cancel()
		return server.Shutdown(shutdownCtx)
	})

	address := "0.0.0.0:" + config.Server.Port
	g.Go(func() error {
		slog.Info("Starting server at " + address)
		if err := server.Start(address); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})

	if err := g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		slog.Error("Server exited with error", "err", err)
		return err
	}

	slog.Info("Server exited gracefully")
	return nil
}
