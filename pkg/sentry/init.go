package sentry

import (
	"go-echo-template/pkg/environment"

	"github.com/getsentry/sentry-go"
)

func Init(dsn string, environment environment.Type) error {
	if dsn == "" {
		return nil
	}
	tracesSampleRate := 0.7
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		TracesSampleRate: tracesSampleRate,
		AttachStacktrace: true,
		Environment:      string(environment),
	})
	if err != nil {
		return err
	}
	return nil
}
