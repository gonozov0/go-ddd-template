package products

import (
	"context"

	"fmt"

	domain "go-ddd-template/internal/domain/products"
	"go-ddd-template/internal/domain/shared/valueobjects"
)

type ProductToCreate struct {
	Name  valueobjects.ProductName
	Price valueobjects.ProductPrice
}

type ProductsToCreate []ProductToCreate

func (s ProductService) CreateProducts(
	ctx context.Context,
	productsToCreate ProductsToCreate,
) (valueobjects.ProductIDs, error) {
	var productIds valueobjects.ProductIDs

	err := s.repo.CreateProducts(ctx, func() (domain.Products, error) {
		products := make([]domain.Product, 0, len(productsToCreate))

		for _, productToCreate := range productsToCreate {
			product := domain.CreateProduct(productToCreate.Name, productToCreate.Price)
			products = append(products, *product)
			productIds = append(productIds, product.GetID())
		}

		return products, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create products: %w", err)
	}

	return productIds, nil
}
