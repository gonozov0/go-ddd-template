package sentry

import (
	"log/slog"

	"github.com/getsentry/sentry-go"

	"fmt"

	"go-ddd-template/pkg/environment"
)

const (
	language = "go"
)

// Init initializes sentry client.
func Init(cfg Config) error {
	err := sentry.Init(
		sentry.ClientOptions{
			Dsn:              cfg.DSN,
			Environment:      getSentryEnvironment(cfg.Environment),
			AttachStacktrace: true,
			Release:          cfg.Release,
			BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
				event.Tags["project"] = cfg.Project
				event.Tags["language"] = language
				event.Tags["service"] = cfg.Service

				return event
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to init sentry: %w", err)
	}

	slog.Info(
		"error booster initialized",
		slog.String("project", cfg.Project),
		slog.String("service", cfg.Service),
		slog.String("release", cfg.Release),
		slog.String("environment", string(cfg.Environment)),
	)

	return nil
}

// getSentryEnvironment returns sentry environment name for given environment type.
// available environments: development|testing|prestable|production|pre-production
func getSentryEnvironment(envType environment.Type) string {
	switch envType {
	case environment.Production:
		return "production"
	case environment.Testing:
		return "testing"
	case environment.Dev:
		return "development"
	default:
		return ""
	}
}
