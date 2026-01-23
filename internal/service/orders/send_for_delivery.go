package orders

import (
	"context"

	"fmt"

	domain "go-ddd-template/internal/domain/orders"
)

func (s OrderService) HandleCreatedOrders() error {
	if err := s.orderRepo.ProcessOrders(context.Background(), func(orders []*domain.Order) error {
		for _, order := range orders {
			if err := order.Process(); err != nil {
				return fmt.Errorf("failed to set processing status: %w", err)
			}
		}

		return nil
	}); err != nil {
		return fmt.Errorf("failed send order created events: %w", err)
	}

	return nil
}
