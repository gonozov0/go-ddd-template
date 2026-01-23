package deliveries

import (
	"context"

	"fmt"

	domain "go-ddd-template/internal/domain/deliveries"
	"go-ddd-template/internal/domain/shared/valueobjects"
)

type DeliveryToCreate struct {
	OrderID valueobjects.OrderID
}

func (s DeliveryService) CreateDelivery(
	ctx context.Context,
	deliveryToCreate DeliveryToCreate,
) (*domain.Delivery, error) {
	delivery, err := s.repo.CreateDelivery(ctx, func() (*domain.Delivery, error) {
		delivery := domain.CreateDelivery(deliveryToCreate.OrderID)

		return delivery, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create delivery: %w", err)
	}

	return delivery, nil
}
