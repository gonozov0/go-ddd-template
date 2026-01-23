package crontests

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"go-ddd-template/internal/application"
	orders "go-ddd-template/internal/application/orders/crons"
)

type OrdersSuite struct {
	suite.Suite
	application.ServerSuite

	CronHandlers orders.OrderHandlers
}

func (s *OrdersSuite) SetupTest() {
	s.ServerSuite.SetupTest()
	s.CronHandlers = orders.SetupHandlers(s.OrdersRepo, s.UsersRepo, 0)
}

func TestOrdersSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(OrdersSuite))
}
