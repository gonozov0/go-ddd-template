package deliveries

import (
	"context"

	domain "go-ddd-template/internal/domain/deliveries"
	"go-ddd-template/internal/domain/shared/valueobjects"
)

type DeliveryRepository interface {
	CreateDelivery(
		ctx context.Context,
		createFn func() (*domain.Delivery, error),
	) (*domain.Delivery, error)
	ListDeliveries(ctx context.Context) ([]domain.Delivery, error)
	DeleteDelivery(ctx context.Context, deliveryID valueobjects.DeliveryID) error
}

type DeliveryService struct {
	repo DeliveryRepository
}

func NewDeliveryService(r DeliveryRepository) DeliveryService {
	return DeliveryService{
		repo: r,
	}
}
