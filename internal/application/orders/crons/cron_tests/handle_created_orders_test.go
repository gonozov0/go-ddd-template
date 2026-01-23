package crontests

import (
	"context"

	orderunithelpers "go-ddd-template/internal/application/orders/shared/helpers"
	userunithelpers "go-ddd-template/internal/application/users/shared/helpers"
	domain "go-ddd-template/internal/domain/orders"
	"go-ddd-template/internal/domain/shared/valueobjects"
)

func (s *OrdersSuite) TestHandleCreatedOrders() {
	order := s.createOrder()

	s.Run("With valid state should update order status to processing", func() {
		err := s.CronHandlers.HandleCreatedOrders()
		s.Require().NoError(err)

		actualOrder, err := s.OrdersRepo.GetOrder(context.Background(), order.GetID())
		s.Require().NoError(err)

		s.Require().Equal(domain.OrderStatusProcessing, actualOrder.GetStatus())
	})
}

func (s *OrdersSuite) createOrder() *domain.Order {
	user := userunithelpers.CreateUser(s, s.UsersRepo)
	products := orderunithelpers.CreateProducts(s, s.OrdersRepo,
		[]orderunithelpers.GenerateProductOption{
			orderunithelpers.ProductWithStatus(valueobjects.ProductStatusPublished),
		},
		[]orderunithelpers.GenerateProductOption{
			orderunithelpers.ProductWithStatus(valueobjects.ProductStatusPublished),
		},
	)
	order := orderunithelpers.GenerateOrder(s,
		orderunithelpers.OrderWithUserID(user.GetID()),
		orderunithelpers.OrderWithProducts(products),
	)
	_, err := s.OrdersRepo.CreateOrder(
		context.Background(),
		order.GetProductIDs(),
		func(p []domain.Product) (*domain.Order, error) { return order, nil },
	)
	s.Require().NoError(err)

	return order
}
