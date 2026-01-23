package helpers

import (
	"context"

	domain "go-ddd-template/internal/domain/products"
	"go-ddd-template/pkg/testify"
)

type ProductsCreater interface {
	CreateProducts(
		ctx context.Context,
		createFn func() (domain.Products, error),
	) error
}

// CreateRandomProducts - создает count случайных продуктов
func CreateRandomProducts(s testify.Suite, repo ProductsCreater, count int) domain.Products {
	products := make(domain.Products, 0, count)
	for range count {
		products = append(products, *GenerateProduct(s))
	}

	s.Require().
		NoError(repo.CreateProducts(context.Background(), func() (domain.Products, error) { return products, nil }))

	return products
}
