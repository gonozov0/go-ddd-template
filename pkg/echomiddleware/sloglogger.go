package echomiddleware

import (
	"context"
	"log/slog"
	"net/http"

	"go-echo-ddd-template/pkg/contextkeys"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	parentSeparatorNumber = 3 // https://www.w3.org/TR/trace-context/#version-format
)

type logger interface {
	LogAttrs(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr)
}

func SlogLoggerMiddleware(logger logger) echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:    true,
		LogURI:       true,
		LogMethod:    true,
		LogRemoteIP:  true,
		LogProtocol:  true,
		LogUserAgent: true,
		LogLatency:   true,
		LogError:     true,
		LogHeaders:   []string{AwsRequestIDHeader, RequestIDHeader, TraceParentHeader},
		HandleError:  true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			attrs := []any{
				slog.String("path", v.URI),
				slog.Int("status_code", v.Status),
				slog.String("method", v.Method),
				slog.String("protocol", v.Protocol),
				slog.String("remote_ip", v.RemoteIP),
				slog.String("user_agent", v.UserAgent),
				slog.String("exec_time", v.Latency.String()),
			}
			level := slog.LevelInfo
			msg := "REQUEST"

			// Adding request_id and trace_id
			reqID := getRequestID(v.Headers)
			attrs = append(attrs, slog.String(string(contextkeys.RequestIDCtxKey), reqID))
			traceID := getTraceID(v.Headers)
			attrs = append(attrs, slog.String(string(contextkeys.TraceIDCtxKey), traceID))

			// Adding container_id attribute
			if containerID := c.Param("containerID"); containerID != "" {
				attrs = append(attrs, slog.String("container_id", containerID))
			}

			respErrStr := "?"
			if v.Error != nil {
				respErrStr = v.Error.Error()
				attrs = append(attrs, slog.String("err", respErrStr))
			}

			// Change level on 5xx
			if v.Status >= http.StatusInternalServerError {
				level = slog.LevelError
				msg = "REQUEST_ERROR: " + respErrStr
			}
			logger.LogAttrs(c.Request().Context(), level, msg, slog.Group("context", attrs...))
			return nil
		},
	})
}
