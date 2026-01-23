package deliveries

import (
	"context"

	"go-ddd-template/internal/domain/shared/valueobjects"
)

func (s DeliveryService) DeleteDelivery(ctx context.Context, deliveryID valueobjects.DeliveryID) error {
	return s.repo.DeleteDelivery(ctx, deliveryID)
}
