package orders_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"go-ddd-template/internal/application"
	orders "go-ddd-template/internal/application/orders/server"
)

type OrdersSuite struct {
	suite.Suite
	application.ServerSuite

	GRPCHandlers orders.OrderHandlers
}

func (s *OrdersSuite) SetupTest() {
	s.ServerSuite.SetupTest()
	s.GRPCHandlers = orders.SetupHandlers(s.OrdersRepo, s.UsersRepo)
}

func TestOrdersSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(OrdersSuite))
}
