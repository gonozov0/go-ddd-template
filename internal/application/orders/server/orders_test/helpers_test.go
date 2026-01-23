package orders_test

import (
	"context"

	orderunithelpers "go-ddd-template/internal/application/orders/shared/helpers"
	userunithelpers "go-ddd-template/internal/application/users/shared/helpers"
	domain "go-ddd-template/internal/domain/orders"
	"go-ddd-template/internal/domain/shared/valueobjects"
)

func (s *OrdersSuite) prepareOrder() *domain.Order {
	user := userunithelpers.CreateUser(s, s.UsersRepo, userunithelpers.UserWithID(s.UserID))
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

	return order
}

func (s *OrdersSuite) createOrder() *domain.Order {
	order := s.prepareOrder()
	_, err := s.OrdersRepo.CreateOrder(
		context.Background(),
		order.GetProductIDs(),
		func(p []domain.Product) (*domain.Order, error) { return order, nil },
	)
	s.Require().NoError(err)

	return order
}
