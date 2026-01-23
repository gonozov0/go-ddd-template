package loggerutils

import (
	"log/slog"
)

type AttrError struct {
	err   error
	attrs []slog.Attr
}

func NewAttrError(err error, attrs ...slog.Attr) AttrError {
	return AttrError{
		err:   err,
		attrs: attrs,
	}
}

func (e AttrError) Error() string {
	if e.err == nil {
		return ""
	}

	return e.err.Error()
}

func (e AttrError) Unwrap() error {
	return e.err
}

func (e AttrError) Attrs() []slog.Attr {
	if e.err == nil {
		return nil
	}

	return e.attrs
}
