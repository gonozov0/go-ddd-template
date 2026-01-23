package helpers

import (
	"context"

	domain "go-ddd-template/internal/domain/orders"
	"go-ddd-template/internal/domain/shared/valueobjects"
	"go-ddd-template/pkg/testify"
)

type ProductsCreater interface {
	CreateProducts(
		ctx context.Context,
		createFn func() (domain.Products, error),
	) error
}

func CreateProducts(s testify.Suite, repo ProductsCreater, optss ...[]GenerateProductOption) domain.Products {
	products := make(domain.Products, 0, len(optss))
	for _, opts := range optss {
		products = append(products, *GenerateProduct(s, opts...))
	}

	s.Require().
		NoError(repo.CreateProducts(context.Background(), func() (domain.Products, error) { return products, nil }))

	return products
}

type OrderCreater interface {
	CreateOrder(
		ctx context.Context,
		productIDs valueobjects.ProductIDs,
		createFn func([]domain.Product) (*domain.Order, error),
	) (*domain.Order, error)
}

func CreateOrder(s testify.Suite, repo OrderCreater, opts ...GenerateOrderOption) domain.Order {
	order := GenerateOrder(s, opts...)

	_, err := repo.CreateOrder(
		context.Background(),
		order.GetProductIDs(),
		func(p []domain.Product) (*domain.Order, error) { return order, nil },
	)
	s.Require().NoError(err)

	return *order
}
