package deliveries

import (
	userunithelpers "go-ddd-template/internal/application/users/shared/helpers"
)

func (s *DeliverySuite) TestDeliveryCreation() {
	s.userHelper.CreateUser(userunithelpers.UserWithID(s.UserID))

	productIDs := s.productHelper.CreateRandomProducts(3)

	orderID := s.orderHelper.CreateOrder(s.UserCtx, productIDs)

	s.Run("Wait For Delivery Creation", func() {
		deliveryID := s.deliveryHelper.WaitForDeliveryCreation(orderID)
		s.Require().NotEmpty(deliveryID)
	})
}
