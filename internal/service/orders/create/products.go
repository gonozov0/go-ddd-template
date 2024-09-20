package create

import (
	"context"
	"fmt"

	"go-echo-template/internal/domain/products"

	"github.com/google/uuid"
)

func (s *OrderCreationService) reserveProducts(ctx context.Context, items []Item) ([]products.Product, error) {
	ids := getItemIDs(items)
	ps, err := s.productRepo.GetProducts(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}

	if err := s.orderRepo.ReserveProducts(ctx, ids); err != nil {
		return nil, fmt.Errorf("failed to reserve products: %w", err)
	}

	return ps, nil
}

func getItemIDs(items []Item) []uuid.UUID {
	itemIDs := make([]uuid.UUID, 0, len(items))
	for _, i := range items {
		itemIDs = append(itemIDs, i.ID)
	}
	return itemIDs
}
