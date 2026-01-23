package crons

import (
	"fmt"
)

func (h *OrderHandlers) HandleCreatedOrders() error {
	if err := h.orderService.HandleCreatedOrders(); err != nil {
		return fmt.Errorf("failed send order created events: %w", err)
	}

	return nil
}
