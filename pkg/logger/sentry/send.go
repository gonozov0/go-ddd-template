package sentry

import (
	"errors"
	"log/slog"
	"reflect"
	"runtime"

	"github.com/getsentry/sentry-go"

	loggerutils "go-ddd-template/pkg/logger/utils"
)

type Event struct {
	Message string
	Err     error
	Attrs   map[string]any
	Level   slog.Level
}

func Send(event Event) {
	sentryEvent := sentry.NewEvent()
	sentryEvent.Level = getSentryLevel(event.Level)
	sentryEvent.Message = event.Message

	if event.Err != nil {
		sentryEvent.Extra["err"] = event.Err.Error()
		sentryEvent.Exception = []sentry.Exception{
			{
				Type:       reflect.TypeOf(event.Err).String(),
				Value:      event.Err.Error(),
				Stacktrace: makeSentryStacktrace(event.Err),
			},
		}
	}

	for key, value := range event.Attrs {
		if key == loggerutils.KeyErr {
			continue
		}

		sentryEvent.Extra[key] = value
	}

	sentry.CaptureEvent(sentryEvent)
}

func makeSentryStacktrace(err error) *sentry.Stacktrace {
	var runtimeErr runtime.Error

	if errors.As(err, &runtimeErr) {
		return sentry.NewStacktrace()
	}

	return nil
}

func getSentryLevel(level slog.Level) sentry.Level {
	switch level {
	case slog.LevelDebug:
		return sentry.LevelDebug
	case slog.LevelInfo:
		return sentry.LevelInfo
	case slog.LevelWarn:
		return sentry.LevelWarning
	case slog.LevelError:
		return sentry.LevelError
	default:
		return sentry.LevelFatal
	}
}
