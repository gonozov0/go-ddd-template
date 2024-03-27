package sentry

import (
	"log/slog"

	"github.com/getsentry/sentry-go"
)

func Init(dsn, environment string) {
	if dsn == "" {
		return
	}
	tracesSampleRate := 0.7
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		TracesSampleRate: tracesSampleRate,
		AttachStacktrace: true,
		Environment:      environment,
	})
	if err != nil {
		slog.Error("sentry.Init: %s", err)
	}
}
