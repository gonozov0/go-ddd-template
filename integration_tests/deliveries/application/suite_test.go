package deliveries

import (
	"testing"

	"github.com/stretchr/testify/suite"

	deliveryhelper "go-ddd-template/integration_tests/deliveries/helpers"
	orderhelper "go-ddd-template/integration_tests/orders/helpers"
	producthelper "go-ddd-template/integration_tests/products/helpers"
	"go-ddd-template/integration_tests/suites"
	userhelper "go-ddd-template/integration_tests/users/helpers"
)

type DeliverySuite struct {
	suites.BaseSuite

	userHelper     *userhelper.HelperSuite
	productHelper  *producthelper.HelperSuite
	orderHelper    *orderhelper.HelperSuite
	deliveryHelper *deliveryhelper.HelperSuite
}

func (s *DeliverySuite) SetupSuite() {
	s.BaseSuite.SetupSuite()

	s.userHelper = userhelper.NewHelperSuite(&s.BaseSuite)
	s.productHelper = producthelper.NewHelperSuite(&s.BaseSuite)
	s.orderHelper = orderhelper.NewHelperSuite(&s.BaseSuite)
	s.deliveryHelper = deliveryhelper.NewHelperSuite(&s.BaseSuite)
}

func TestConsumerSuite(t *testing.T) {
	suite.Run(t, new(DeliverySuite))
}
