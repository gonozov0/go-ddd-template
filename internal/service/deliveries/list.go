package deliveries

import (
	"context"

	"fmt"

	"go-ddd-template/internal/domain/deliveries"
)

func (s DeliveryService) ListDeliveries(ctx context.Context) ([]deliveries.Delivery, error) {
	deliveries, err := s.repo.ListDeliveries(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list deliveries: %w", err)
	}

	return deliveries, nil
}
