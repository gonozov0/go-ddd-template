package server_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"go-ddd-template/internal/application"
	"go-ddd-template/internal/application/deliveries/consumers"
	"go-ddd-template/internal/application/deliveries/server"
)

type DeliveriesSuite struct {
	suite.Suite
	application.ServerSuite

	GRPCHandlers      server.DeliveryHandlers
	DeliveryConsumers consumers.DeliveryConsumers
}

func (s *DeliveriesSuite) SetupTest() {
	s.ServerSuite.SetupTest()
	s.GRPCHandlers = server.SetupHandlers(s.DeliveriesRepo)
	s.DeliveryConsumers = *consumers.NewDeliveryConsumers(s.DeliveriesRepo)
}

func TestOrdersSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(DeliveriesSuite))
}
