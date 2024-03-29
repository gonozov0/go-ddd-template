package internal

import (
	"log/slog"
	"net/http"

	"go-echo-ddd-template/internal/application/orders"
	"go-echo-ddd-template/internal/application/users"
	ordersInfra "go-echo-ddd-template/internal/infrastructure/orders"
	productsInfra "go-echo-ddd-template/internal/infrastructure/products"
	usersInfra "go-echo-ddd-template/internal/infrastructure/users"
	"go-echo-ddd-template/pkg/echomiddleware"

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

	userRepo := usersInfra.NewPostgresRepo()
	users.Setup(e, userRepo)

	productRepo := productsInfra.NewPostgresRepo()
	orderRepo := ordersInfra.NewPostgresRepo()
	orders.Setup(e, orderRepo, userRepo, productRepo)

	return e
}
