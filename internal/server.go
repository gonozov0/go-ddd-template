package internal

import (
	"log/slog"
	"net/http"

	"go-echo-template/internal/application/users"
	usersInfra "go-echo-template/internal/infrastructure/users"
	"go-echo-template/pkg/echomiddleware"

	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func newServer(_ Config) *echo.Echo {
	e := echo.New()

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(echomiddleware.SlogLoggerMiddleware(slog.Default()))
	e.Use(echomiddleware.PutRequestIDContext)
	e.Use(middleware.Recover())
	e.Use(sentryecho.New(sentryecho.Options{Repanic: true}))
	e.Use(echomiddleware.PutSentryContext)

	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	usersRepo := usersInfra.NewPostgresRepo()
	users.Setup(e, usersRepo)

	return e
}
