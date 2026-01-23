package loggerutils

import "log/slog"

func Attr(key string, value any) slog.Attr {
	return slog.Attr{
		Key:   key,
		Value: slog.AnyValue(value),
	}
}
