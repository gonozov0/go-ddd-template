package products

import (
	"context"

	"fmt"

	domain "go-ddd-template/internal/domain/products"
	"go-ddd-template/internal/domain/shared/valueobjects"
)

func (s ProductService) CheckProduct(
	ctx context.Context,
	productID valueobjects.ProductID,
) error {
	if err := s.repo.UpdateProduct(ctx, productID, func(product *domain.Product) error {
		// Проверка и валидация товара на корректность
		if err := product.Publish(); err != nil {
			return fmt.Errorf("failed to publish product: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	return nil
}
