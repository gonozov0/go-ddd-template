package products

import (
	"context"

	"fmt"

	domain "go-ddd-template/internal/domain/products"
	"go-ddd-template/internal/domain/shared/valueobjects"
)

func (s ProductService) DeleteProducts(ctx context.Context, ids valueobjects.ProductIDs) error {
	if err := s.repo.DeleteProducts(ctx, ids, func(products []domain.Product) error {
		for _, product := range products {
			if product.GetImageFilename().IsEmpty() || product.GetImageURL().IsEmpty() {
				continue
			}

			if err := s.imageStorage.DeleteImage(ctx, product.GetImageFilename()); err != nil {
				return fmt.Errorf("failed to delete image: %w", err)
			}
		}

		return nil
	}); err != nil {
		return fmt.Errorf("failed to delete products: %w", err)
	}

	return nil
}
