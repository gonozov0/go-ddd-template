package application

import (
	ordersInfra "go-echo-ddd-template/internal/infrastructure/orders"
	productsInfra "go-echo-ddd-template/internal/infrastructure/products"
	usersInfra "go-echo-ddd-template/internal/infrastructure/users"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
)

type ServerSuite struct {
	suite.Suite
	HTTPServer   *echo.Echo
	GRPCServer   *grpc.Server
	UsersRepo    *usersInfra.InMemoryRepo
	OrdersRepo   *ordersInfra.InMemoryRepo
	ProductsRepo *productsInfra.InMemoryRepo
}

func (s *ServerSuite) SetupSuite() {
	s.UsersRepo = usersInfra.NewInMemoryRepo()
	s.OrdersRepo = ordersInfra.NewInMemoryRepo()
	s.ProductsRepo = productsInfra.NewInMemoryRepo()
	s.HTTPServer = SetupHTTPServer(s.UsersRepo, s.OrdersRepo, s.ProductsRepo)
	s.GRPCServer = SetupGRPCServer(s.UsersRepo, s.OrdersRepo, s.ProductsRepo)
}
