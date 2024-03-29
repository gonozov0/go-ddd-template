package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/getsentry/sentry-go"

	"go-echo-ddd-template/pkg/contextkeys"
)

func Setup() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	logger := slog.New(newSentryJSONCtxHandler(os.Stdout, opts))
	slog.SetDefault(logger)
}

type sentryJSONCtxHandler struct {
	*slog.JSONHandler
}

func (h *sentryJSONCtxHandler) Handle(ctx context.Context, r slog.Record) error {
	contextInRecord := false
	r.Attrs(func(a slog.Attr) bool {
		contextInRecord = a.Key == "context"
		return !contextInRecord
	})

	requestIDKey := contextkeys.RequestIDCtxKey
	traceIDKey := contextkeys.TraceIDCtxKey
	requestID := fmt.Sprintf("%v", ctx.Value(requestIDKey))
	traceID := fmt.Sprintf("%v", ctx.Value(traceIDKey))

	// Sending event to sentry
	if r.Level == slog.LevelError {
		if hub, ok := ctx.Value(sentry.HubContextKey).(*sentry.Hub); ok && hub != nil {
			hub.WithScope(func(scope *sentry.Scope) {
				scope.SetLevel(sentry.LevelError)
				scope.SetTag(string(requestIDKey), requestID)
				scope.SetTag(string(traceIDKey), traceID)
				hub.CaptureMessage(r.Message)
			})
		}
	}

	// Log recored does not contain any custom 'context' key -> logging values from ctx
	if !contextInRecord {
		r.AddAttrs(slog.Group(
			"context",
			slog.String(string(requestIDKey), requestID),
			slog.String(string(traceIDKey), traceID),
		))
	}

	return h.JSONHandler.Handle(ctx, r)
}

func newSentryJSONCtxHandler(w io.Writer, opts *slog.HandlerOptions) *sentryJSONCtxHandler {
	jsonHandler := slog.NewJSONHandler(w, opts)
	return &sentryJSONCtxHandler{
		JSONHandler: jsonHandler,
	}
}
