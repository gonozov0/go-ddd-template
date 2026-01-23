package orders

import (
	"context"

	"fmt"

	"go-ddd-template/internal/domain/shared/valueobjects"
)

func (s OrderService) DeleteOrder(
	ctx context.Context,
	orderID valueobjects.OrderID,
) error {
	err := s.orderRepo.DeleteOrder(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}

	return nil
}
