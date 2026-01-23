package products

import (
	"context"

	"fmt"

	"go-ddd-template/internal/domain/products"
	"go-ddd-template/internal/domain/shared/valueobjects"
)

func (s ProductService) GetProducts(ctx context.Context, ids valueobjects.ProductIDs) ([]products.Product, error) {
	products, err := s.repo.GetProducts(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}

	return products, nil
}
