package sentry

import (
	"github.com/getsentry/sentry-go"
)

func Init(dsn, environment string) error {
	if dsn == "" {
		return nil
	}
	tracesSampleRate := 0.7
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		TracesSampleRate: tracesSampleRate,
		AttachStacktrace: true,
		Environment:      environment,
	})
	if err != nil {
		return err
	}
	return nil
}
