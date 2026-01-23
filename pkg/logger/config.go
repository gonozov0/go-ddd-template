package logger

import (
	"log/slog"
	"strings"

	"go-ddd-template/pkg/logger/sentry"
)

type Config struct {
	EnableJSONFormat bool
	LogLevel         slog.Level
	Sentry           sentry.Config
}

func GetLogLevel(lvl string) slog.Level {
	switch strings.ToLower(lvl) {
	case "debug":
		return slog.LevelDebug
	case "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
