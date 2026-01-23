package orders

import (
	"context"

	"fmt"

	domain "go-ddd-template/internal/domain/orders"
	"go-ddd-template/internal/domain/shared/valueobjects"
)

type OrderToCreate struct {
	UserID     valueobjects.UserID
	ProductIDs valueobjects.ProductIDs
}

func (s OrderService) CreateOrder(ctx context.Context, orderToCreate OrderToCreate) (*domain.Order, error) {
	// check if user exists
	_, err := s.userService.GetUser(ctx, orderToCreate.UserID)
	if err != nil {
		return nil, err
	}

	order, err := s.orderRepo.CreateOrder(
		ctx,
		orderToCreate.ProductIDs,
		func(products []domain.Product) (*domain.Order, error) {
			for i, product := range products {
				if err := product.Reserve(); err != nil {
					return nil, fmt.Errorf("failed to reserve product: %w", err)
				}

				products[i] = product
			}

			return domain.CreateOrder(orderToCreate.UserID, products), nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return order, nil
}
