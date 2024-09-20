package orders

import (
	"context"

	"go-echo-template/internal/domain/orders"

	"github.com/google/uuid"
)

type order struct {
	userID uuid.UUID
	status orders.OrderStatus
	items  []orders.Item
}

type InMemoryRepo struct {
	orders map[uuid.UUID]order
}

func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{
		orders: make(map[uuid.UUID]order),
	}
}

func (r *InMemoryRepo) SaveOrder(_ context.Context, o *orders.Order) error {
	r.orders[o.ID()] = order{
		userID: o.UserID(),
		status: o.Status(),
		items:  o.Items(),
	}

	return nil
}

func (r *InMemoryRepo) GetOrder(_ context.Context, id uuid.UUID) (*orders.Order, error) {
	o, ok := r.orders[id]
	if !ok {
		return nil, orders.ErrOrderNotFound
	}

	order, err := orders.NewOrder(id, o.userID, o.status, o.items)
	if err != nil {
		return nil, err
	}

	return order, nil
}
