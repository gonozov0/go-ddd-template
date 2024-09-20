package orders

import (
	"context"
	"fmt"

	"go-echo-template/internal/domain/orders"

	"github.com/google/uuid"
)

type product struct {
	Available bool
}

type order struct {
	UserID uuid.UUID
	Status orders.OrderStatus
	Items  []orders.Item
}

type InMemoryRepo struct {
	orders   map[uuid.UUID]order
	products map[uuid.UUID]product
}

func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{
		orders:   make(map[uuid.UUID]order),
		products: make(map[uuid.UUID]product),
	}
}

func (r *InMemoryRepo) SaveOrder(_ context.Context, o *orders.Order) error {
	r.orders[o.ID()] = order{
		UserID: o.UserID(),
		Status: o.Status(),
		Items:  o.Items(),
	}

	return nil
}

func (r *InMemoryRepo) GetOrder(_ context.Context, id uuid.UUID) (*orders.Order, error) {
	o, ok := r.orders[id]
	if !ok {
		return nil, orders.ErrOrderNotFound
	}

	order, err := orders.NewOrder(id, o.UserID, o.Status, o.Items)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (r *InMemoryRepo) ReserveProducts(_ context.Context, ids []uuid.UUID) error {
	for _, id := range ids {
		p, ok := r.products[id]
		if !ok {
			return fmt.Errorf("%w: id %s", orders.ErrProductNotFound, id)
		}
		if !p.Available {
			return fmt.Errorf("%w: id %s", orders.ErrProductAlreadyReserved, id)
		}
		r.products[id] = product{
			Available: false,
		}
	}

	return nil
}
