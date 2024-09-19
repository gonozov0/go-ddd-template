package application

import (
	"log/slog"
	"net/http"

	"go-echo-template/generated/openapi"
	"go-echo-template/internal/application/orders"
	"go-echo-template/internal/application/users"
	"go-echo-template/pkg/echomiddleware"

	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type httpServer struct {
	users.UserHandlers
	orders.OrderHandlers
}

func SetupHTTPServer(userRepo UserRepository, orderRepo OrderRepository, productRepo ProductRepository) *echo.Echo {
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

	server := httpServer{}
	server.UserHandlers = users.SetupHandlers(userRepo)
	server.OrderHandlers = orders.SetupHandlers(orderRepo, userRepo, productRepo)

	openapi.RegisterHandlers(e, server)

	return e
}
