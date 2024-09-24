package orders

import (
	"context"
	"fmt"

	ordersDomain "go-echo-template/internal/domain/orders"
	productsDomain "go-echo-template/internal/domain/products"

	"github.com/google/uuid"
)

type Item struct {
	ID uuid.UUID
}

func (s *OrderCreationService) CreateOrder(
	ctx context.Context,
	userID uuid.UUID,
	items []Item,
) (*ordersDomain.Order, error) {
	// check if user exists
	_, err := s.userRepo.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	itemIDs := getItemIDs(items)
	order, err := s.orderRepo.CreateOrder(ctx, itemIDs, func() (*ordersDomain.Order, error) {
		products, err := s.productRepo.GetProducts(ctx, itemIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to get products: %w", err)
		}

		orderItems, err := makeOrderItems(items, products)
		if err != nil {
			return nil, fmt.Errorf("failed to create order items: %w", err)
		}
		return ordersDomain.CreateOrder(userID, orderItems)
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return order, nil
}

func makeOrderItems(items []Item, ps []productsDomain.Product) ([]ordersDomain.Item, error) {
	orderItems := make([]ordersDomain.Item, 0, len(items))
	for i, product := range ps {
		item := items[i]
		orderItem, err := ordersDomain.NewItem(item.ID, product.Name(), product.Price())
		if err != nil {
			return nil, err
		}
		orderItems = append(orderItems, *orderItem)
	}
	return orderItems, nil
}

func getItemIDs(items []Item) []uuid.UUID {
	itemIDs := make([]uuid.UUID, 0, len(items))
	for _, i := range items {
		itemIDs = append(itemIDs, i.ID)
	}
	return itemIDs
}
