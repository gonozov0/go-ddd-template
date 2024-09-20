package create

import (
	"context"

	"go-echo-template/internal/domain/orders"
	"go-echo-template/internal/domain/products"

	"github.com/google/uuid"
)

type Item struct {
	ID uuid.UUID
}

func (s *OrderCreationService) CreateOrder(ctx context.Context, userID uuid.UUID, items []Item) (*orders.Order, error) {
	// check if user exists
	_, err := s.userRepo.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	ps, err := s.reserveProducts(ctx, items)
	if err != nil {
		return nil, err
	}

	orderItems, err := makeOrderItems(items, ps)
	if err != nil {
		return nil, err
	}
	order, err := orders.CreateOrder(userID, orderItems)
	if err != nil {
		return nil, err
	}
	if err = s.orderRepo.SaveOrder(ctx, order); err != nil {
		return nil, err
	}

	return order, nil
}

func makeOrderItems(items []Item, ps []products.Product) ([]orders.Item, error) {
	orderItems := make([]orders.Item, 0, len(items))
	for i, product := range ps {
		item := items[i]
		orderItem, err := orders.NewItem(item.ID, product.Name(), product.Price())
		if err != nil {
			return nil, err
		}
		orderItems = append(orderItems, *orderItem)
	}
	return orderItems, nil
}
