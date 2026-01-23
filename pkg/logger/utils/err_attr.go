package loggerutils

import "log/slog"

const KeyErr = "err"

func ErrAttr(err any) slog.Attr {
	return slog.Attr{
		Key:   KeyErr,
		Value: slog.AnyValue(err),
	}
}
