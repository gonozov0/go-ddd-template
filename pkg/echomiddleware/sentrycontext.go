package echomiddleware

import (
	"context"

	"github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
)

// PutSentryContext adding sentryHub to request context.Context.
func PutSentryContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if hub := sentryecho.GetHubFromContext(c); hub != nil {
			hub.ConfigureScope(func(scope *sentry.Scope) {
				scope.SetRequest(c.Request())
			})
			ctx := c.Request().Context()
			ctx = context.WithValue(ctx, sentry.HubContextKey, hub)
			c.SetRequest(c.Request().WithContext(ctx))
		}
		return next(c)
	}
}
