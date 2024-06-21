package application

import (
	ordersInfra "go-echo-ddd-template/internal/infrastructure/orders"
	productsInfra "go-echo-ddd-template/internal/infrastructure/products"
	usersInfra "go-echo-ddd-template/internal/infrastructure/users"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

type ServerSuite struct {
	suite.Suite
	Echo         *echo.Echo
	UsersRepo    *usersInfra.InMemoryRepo
	OrdersRepo   *ordersInfra.InMemoryRepo
	ProductsRepo *productsInfra.InMemoryRepo
}

func (s *ServerSuite) SetupSuite() {
	s.UsersRepo = usersInfra.NewInMemoryRepo()
	s.OrdersRepo = ordersInfra.NewInMemoryRepo()
	s.ProductsRepo = productsInfra.NewInMemoryRepo()
	s.Echo = SetupServer(s.UsersRepo, s.OrdersRepo, s.ProductsRepo)
}
