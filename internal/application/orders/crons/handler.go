package crons

import (
	"context"
	"log/slog"
	"time"

	"golang.org/x/sync/errgroup"

	service "go-ddd-template/internal/service/orders"
	"go-ddd-template/internal/service/users"
	loggerutils "go-ddd-template/pkg/logger/utils"
)

type OrderHandlers struct {
	orderService service.OrderService

	handleCreatedOrdresInterval time.Duration
}

func SetupHandlers(
	ordersRepo service.OrderRepository,
	userRepo users.UserRepository,
	handleCreatedOrdresInterval time.Duration,
) OrderHandlers {
	userService := users.NewUserService(userRepo)
	orderService := service.NewOrderService(ordersRepo, userService)

	return OrderHandlers{
		orderService:                orderService,
		handleCreatedOrdresInterval: handleCreatedOrdresInterval,
	}
}

func (h *OrderHandlers) Start(ctx context.Context, eg *errgroup.Group) {
	eg.Go(func() error {
		ticker := time.NewTicker(h.handleCreatedOrdresInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return nil
			case <-ticker.C:
				err := h.HandleCreatedOrders()
				if err != nil {
					slog.Error("failed to handle created orders", loggerutils.ErrAttr(err))
				} else {
					slog.Info("successfully handled created orders")
				}
			}
		}
	})
}
