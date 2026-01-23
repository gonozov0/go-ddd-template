package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"runtime"
	"runtime/debug"

	"errors"
	"fmt"

	"go-ddd-template/pkg/auth"
	"go-ddd-template/pkg/grpcutils/middlewares"
	"go-ddd-template/pkg/logger/sentry"
	loggerutils "go-ddd-template/pkg/logger/utils"
	"go-ddd-template/pkg/traces"
)

func Setup(cfg Config) error {
	if !cfg.EnableJSONFormat {
		slog.SetLogLoggerLevel(cfg.LogLevel)
		return nil
	}

	if err := sentry.Init(cfg.Sentry); err != nil {
		return fmt.Errorf("failed to init errorbooster: %w", err)
	}

	opts := &slog.HandlerOptions{
		Level: cfg.LogLevel,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.LevelKey:
				a.Key = "levelStr"
				return a
			case slog.TimeKey:
				a.Key = "@timestamp"
				val := a.Value.Time()
				a.Value = slog.Int64Value(val.Unix())

				return a
			default:
				return a
			}
		},
	}
	logger := slog.New(newJSONCtxHandler(os.Stdout, opts))
	slog.SetDefault(logger)

	return nil
}

type JSONCtxHandler struct {
	*slog.JSONHandler
}

//nolint:cyclop
func (h *JSONCtxHandler) Handle(ctx context.Context, r slog.Record) error {
	userInfo := auth.GetUserInfo(ctx)
	if !userInfo.IsEmpty() {
		r.Add("user_tvm_id", userInfo.ID)
	}

	requestID := traces.GetRequestID(ctx)
	r.Add("request_id", requestID)

	correlationID := traces.GetCorrelationID(ctx)
	r.Add("correlation_id", correlationID)

	grpcReqStr := middlewares.GetGRPCRequestStr(ctx)
	r.Add("grpc_request_str", grpcReqStr)

	grpcMethod := middlewares.GetGRPCMethod(ctx)
	r.Add("grpc_method", grpcMethod)

	spanID := traces.GetSpanID(ctx)
	r.Add("span.id", spanID)

	traceID := traces.GetTraceID(ctx)
	r.Add("trace.id", traceID)

	attrs := make(map[string]any, r.NumAttrs())
	r.Attrs(func(attr slog.Attr) bool {
		attrs[attr.Key] = attr.Value.Any()
		return true
	})

	if r.Level == slog.LevelError || r.Level == slog.LevelWarn {
		err := getError(attrs)
		if err != nil {
			var attrErr loggerutils.AttrError

			if errors.As(err, &attrErr) {
				for _, attr := range attrErr.Attrs() {
					if _, exists := attrs[attr.Key]; !exists {
						attrs[attr.Key] = attr.Value.Any()
					}

					r.Add(attr.Key, attr.Value.Any())
				}
			}

			if stacktrace := getStacktrace(err); stacktrace != "" {
				r.Add("stacktrace", stacktrace)
			}
		}

		sentry.Send(sentry.Event{
			Message: r.Message,
			Err:     err,
			Attrs:   attrs,
			Level:   r.Level,
		})
	}

	return h.JSONHandler.Handle(ctx, r)
}

func newJSONCtxHandler(w io.Writer, opts *slog.HandlerOptions) *JSONCtxHandler {
	jsonHandler := slog.NewJSONHandler(w, opts)

	return &JSONCtxHandler{
		JSONHandler: jsonHandler,
	}
}

func getError(attrs map[string]any) error {
	if err, ok := attrs[loggerutils.KeyErr]; ok {
		if e, ok := err.(error); ok {
			return e
		}
	}

	return nil
}

func getStacktrace(err error) string {
	var runtimeErr runtime.Error

	switch {
	case errors.As(err, &runtimeErr):
		return string(debug.Stack())
	default:
		return ""
	}
}
