package orders

import (
	"context"
	"fmt"

	domain "go-echo-template/internal/domain/orders"

	"github.com/google/uuid"
)

type product struct {
	Available bool
}

type order struct {
	UserID uuid.UUID
	Status domain.OrderStatus
	Items  []domain.Item
}

type InMemoryRepo struct {
	domain   map[uuid.UUID]order
	products map[uuid.UUID]product
}

func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{
		domain:   make(map[uuid.UUID]order),
		products: make(map[uuid.UUID]product),
	}
}

func (r *InMemoryRepo) CreateOrder(
	_ context.Context,
	itemIDs []uuid.UUID,
	createFn func() (*domain.Order, error),
) (*domain.Order, error) {
	if err := r.reserveProducts(itemIDs); err != nil {
		return nil, fmt.Errorf("failed to reserve products: %w", err)
	}

	entity, err := createFn()
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	r.domain[entity.ID()] = order{
		UserID: entity.UserID(),
		Status: entity.Status(),
		Items:  entity.Items(),
	}

	return entity, nil
}

func (r *InMemoryRepo) reserveProducts(ids []uuid.UUID) error {
	for _, id := range ids {
		p, ok := r.products[id]
		if !ok {
			return fmt.Errorf("%w: id %s", domain.ErrProductNotFound, id)
		}
		if !p.Available {
			return fmt.Errorf("%w: id %s", domain.ErrProductAlreadyReserved, id)
		}
		r.products[id] = product{
			Available: false,
		}
	}

	return nil
}

func (r *InMemoryRepo) GetOrder(_ context.Context, id uuid.UUID) (*domain.Order, error) {
	o, ok := r.domain[id]
	if !ok {
		return nil, domain.ErrOrderNotFound
	}

	order, err := domain.NewOrder(id, o.UserID, o.Status, o.Items)
	if err != nil {
		return nil, err
	}

	return order, nil
}
