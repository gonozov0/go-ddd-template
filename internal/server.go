package internal

import (
	"log/slog"
	"net/http"

	_ "go-echo-ddd-template/docs" // init swagger
	"go-echo-ddd-template/internal/application/orders"
	"go-echo-ddd-template/internal/application/users"
	ordersInfra "go-echo-ddd-template/internal/infrastructure/orders"
	productsInfra "go-echo-ddd-template/internal/infrastructure/products"
	usersInfra "go-echo-ddd-template/internal/infrastructure/users"
	"go-echo-ddd-template/pkg/echomiddleware"
	"go-echo-ddd-template/pkg/environment"

	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	swagger "github.com/swaggo/echo-swagger"
)

func newServer(config Config) *echo.Echo {
	e := echo.New()

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(echomiddleware.SlogLoggerMiddleware(slog.Default()))
	e.Use(echomiddleware.PutRequestIDContext)
	e.Use(middleware.Recover())
	e.Use(sentryecho.New(sentryecho.Options{Repanic: true}))
	e.Use(echomiddleware.PutSentryContext)

	if config.Server.Environment != environment.Production {
		e.GET("/swagger/*", swagger.WrapHandler)
		e.GET("/swagger", func(c echo.Context) error {
			return c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
		})
	}

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
