package helpers

import (
	"context"

	domain "go-ddd-template/internal/domain/deliveries"
	"go-ddd-template/pkg/testify"
)

type DeliveryCreater interface {
	CreateDelivery(
		ctx context.Context,
		createFn func() (*domain.Delivery, error),
	) (*domain.Delivery, error)
}

func CreateDelivery(s testify.Suite, repo DeliveryCreater, opts ...GenerateDeliveryOption) *domain.Delivery {
	delivery := GenerateDelivery(s, opts...)

	_, err := repo.CreateDelivery(context.Background(), func() (*domain.Delivery, error) { return delivery, nil })
	s.Require().NoError(err)

	return delivery
}
