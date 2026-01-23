package deliveries

import (
	"context"

	"fmt"

	domain "go-ddd-template/internal/domain/deliveries"
	"go-ddd-template/internal/domain/shared/valueobjects"
)

type InMemoryRepo struct {
	deliveries map[valueobjects.DeliveryID]domain.Delivery
}

func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{
		deliveries: make(map[valueobjects.DeliveryID]domain.Delivery),
	}
}

func (r *InMemoryRepo) CreateDelivery(
	ctx context.Context,
	createFn func() (*domain.Delivery, error),
) (*domain.Delivery, error) {
	delivery, err := createFn()
	if err != nil {
		return nil, fmt.Errorf("failed to create delivery: %w", err)
	}

	r.deliveries[delivery.GetID()] = *delivery

	return delivery, nil
}

func (r *InMemoryRepo) ListDeliveries(ctx context.Context) ([]domain.Delivery, error) {
	deliveries := make([]domain.Delivery, 0, len(r.deliveries))

	for _, d := range r.deliveries {
		deliveries = append(deliveries, d)
	}

	return deliveries, nil
}

func (r *InMemoryRepo) DeleteDelivery(ctx context.Context, deliveryID valueobjects.DeliveryID) error {
	delete(r.deliveries, deliveryID)
	return nil
}
