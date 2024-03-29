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

	sentry.Init(config.SentryDSN, config.SentryEnvironment)
	logger.Setup()

	server := newServer(config)
	g, ctx := errgroup.WithContext(context.Background())
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	g.Go(func() error {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), config.InterruptTimeout)
		defer cancel()
		return server.Shutdown(shutdownCtx)
	})

	g.Go(func() error {
		slog.Info("Starting server at 0.0.0.0:8080")
		if err := server.Start("0.0.0.0:8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
