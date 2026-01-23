package consumers

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/google/uuid"

	"errors"
	"fmt"

	"go-ddd-template/internal/domain/shared/valueobjects"
	service "go-ddd-template/internal/service/deliveries"
	loggerutils "go-ddd-template/pkg/logger/utils"
)

type orderCreatedEvent struct {
	OrderID uuid.UUID `json:"order_id"`
}

func (d DeliveryConsumers) HandleOrderProcessing(ctx context.Context, data []byte) error {
	slog.InfoContext(ctx, "HandleOrderProcessing called - message received from topic")

	event, err := d.getEvent(data)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			slog.InfoContext(ctx, "Stop handling order processing events")
			return err
		}

		slog.ErrorContext(ctx,
			"failed to get order_processing event",
			loggerutils.ErrAttr(err),
			loggerutils.Attr("data", string(data)),
		)

		return fmt.Errorf("failed to get order_processing event: %w", err)
	}

	slog.InfoContext(ctx, "Order processing event received", loggerutils.Attr("order_id", event.OrderID.String()))

	orderID, err := valueobjects.NewOrderID(event.OrderID)
	if err != nil {
		slog.ErrorContext(ctx,
			"failed to create order id",
			loggerutils.ErrAttr(err),
			loggerutils.Attr("order_id", event.OrderID.String()),
		)

		return fmt.Errorf("failed to create order id: %w", err)
	}

	delivery, err := d.deliveryService.CreateDelivery(ctx, service.DeliveryToCreate{
		OrderID: orderID,
	})
	if err != nil {
		slog.ErrorContext(ctx,
			"failed to create delivery",
			loggerutils.ErrAttr(err),
			loggerutils.Attr("order_id", event.OrderID.String()),
		)

		return fmt.Errorf("failed to create delivery: %w", err)
	}

	slog.InfoContext(ctx,
		"Delivery created successfully",
		loggerutils.Attr("order_id", event.OrderID.String()),
		loggerutils.Attr("delivery_id", delivery.GetID().String()),
	)

	return nil
}

func (d DeliveryConsumers) getEvent(data []byte) (orderCreatedEvent, error) {
	var event orderCreatedEvent

	if err := json.Unmarshal(data, &event); err != nil {
		return event, fmt.Errorf("failed to unmarshal message: %w", err)
	}

	return event, nil
}
