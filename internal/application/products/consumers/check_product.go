package consumers

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/google/uuid"

	"fmt"

	"go-ddd-template/internal/domain/shared/valueobjects"
)

type productInitedEvent struct {
	ID uuid.UUID `json:"id"`
}

func (d DeliveryConsumers) CheckProduct(ctx context.Context, message []byte) error {
	var event productInitedEvent

	slog.Info("Start process product inited event")

	if err := json.Unmarshal(message, &event); err != nil {
		return fmt.Errorf("failed to unmarshal product inited event: %w", err)
	}

	productID, err := valueobjects.NewProductID(event.ID)
	if err != nil {
		return fmt.Errorf("failed to create product id: %w", err)
	}

	err = d.productService.CheckProduct(ctx, productID)
	if err != nil {
		return fmt.Errorf("failed to check product: %w", err)
	}

	return nil
}
