package create

import (
	"go-echo-ddd-template/internal/domain/orders"
	"go-echo-ddd-template/internal/domain/products"

	"github.com/google/uuid"
)

type Item struct {
	ID uuid.UUID
}

func (s *OrderCreationService) CreateOrder(userID uuid.UUID, items []Item) (*orders.Order, error) {
	// check if user exists
	_, err := s.userRepo.GetUser(userID)
	if err != nil {
		return nil, err
	}

	ps, err := s.reserveProducts(items)
	if err != nil {
		return nil, err
	}
	defer s.productRepo.CancelUpdate()

	orderItems, err := makeOrderItems(items, ps)
	if err != nil {
		return nil, err
	}
	order, err := orders.CreateOrder(userID, orderItems)
	if err != nil {
		return nil, err
	}
	if err = s.orderRepo.SaveOrder(*order); err != nil {
		return nil, err
	}

	if err = s.productRepo.SaveProducts(ps); err != nil {
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
