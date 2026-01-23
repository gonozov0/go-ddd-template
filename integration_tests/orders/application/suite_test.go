package application_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	pb "go-ddd-template/generated/server"
	deliveryhelper "go-ddd-template/integration_tests/deliveries/helpers"
	orderhelper "go-ddd-template/integration_tests/orders/helpers"
	producthelper "go-ddd-template/integration_tests/products/helpers"
	"go-ddd-template/integration_tests/suites"
	userhelper "go-ddd-template/integration_tests/users/helpers"
)

type OrdersSuite struct {
	suites.BaseSuite

	grpcClient pb.OrderServiceClient

	userHelper     *userhelper.HelperSuite
	productHelper  *producthelper.HelperSuite
	deliveryHelper *deliveryhelper.HelperSuite
	orderHelper    *orderhelper.HelperSuite
}

func (s *OrdersSuite) SetupSuite() {
	s.BaseSuite.SetupSuite()

	s.grpcClient = pb.NewOrderServiceClient(s.ServerConn)

	s.userHelper = userhelper.NewHelperSuite(&s.BaseSuite)
	s.productHelper = producthelper.NewHelperSuite(&s.BaseSuite)
	s.deliveryHelper = deliveryhelper.NewHelperSuite(&s.BaseSuite)
	s.orderHelper = orderhelper.NewHelperSuite(&s.BaseSuite)
}

func TestOrdersSuite(t *testing.T) {
	suite.Run(t, new(OrdersSuite))
}
