package application

import (
	"log/slog"
	"net/http"

	"go-echo-ddd-template/generated/openapi"
	"go-echo-ddd-template/internal/application/orders"
	"go-echo-ddd-template/internal/application/users"
	ordersDomain "go-echo-ddd-template/internal/domain/orders"
	productsDomain "go-echo-ddd-template/internal/domain/products"
	usersDomain "go-echo-ddd-template/internal/domain/users"
	"go-echo-ddd-template/pkg/echomiddleware"

	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type UserRepository interface {
	GetUser(id uuid.UUID) (*usersDomain.User, error)
	SaveUser(u usersDomain.User) error
}

type OrderRepository interface {
	SaveOrder(o ordersDomain.Order) error
	GetOrder(id uuid.UUID) (*ordersDomain.Order, error)
}

type ProductRepository interface {
	GetProductsForUpdate(ids []uuid.UUID) ([]productsDomain.Product, error)
	SaveProducts(ps []productsDomain.Product) error
	CancelUpdate()
}

type Server struct {
	users.UserHandlers
	orders.OrderHandlers
}

func SetupServer(userRepo UserRepository, orderRepo OrderRepository, productRepo ProductRepository) *echo.Echo {
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

	server := Server{}
	server.UserHandlers = users.SetupHandlers(userRepo)
	server.OrderHandlers = orders.SetupHandlers(orderRepo, userRepo, productRepo)

	openapi.RegisterHandlers(e, server)

	return e
}
