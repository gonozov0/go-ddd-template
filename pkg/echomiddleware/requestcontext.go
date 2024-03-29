package echomiddleware

import (
	"context"

	"go-echo-ddd-template/pkg/contextkeys"

	"github.com/labstack/echo/v4"
)

func PutRequestIDContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		reqID := getRequestID(c.Request().Header)
		traceID := getTraceID(c.Request().Header)
		ctx := c.Request().Context()
		ctx = context.WithValue(ctx, contextkeys.RequestIDCtxKey, reqID)
		ctx = context.WithValue(ctx, contextkeys.TraceIDCtxKey, traceID)
		c.SetRequest(c.Request().WithContext(ctx))
		return next(c)
	}
}
