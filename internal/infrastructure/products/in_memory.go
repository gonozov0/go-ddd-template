package products

import (
	"context"
	"fmt"

	"go-echo-template/internal/domain/products"

	"github.com/google/uuid"
)

type product struct {
	Name  string
	Price float64
}

type InMemoryRepo struct {
	products map[uuid.UUID]product
}

func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{
		products: make(map[uuid.UUID]product),
	}
}

func (r *InMemoryRepo) CreateProducts(_ context.Context, createFn func() ([]products.Product, error)) error {
	ps, err := createFn()
	if err != nil {
		return fmt.Errorf("failed to create products: %w", err)
	}
	for _, p := range ps {
		r.products[p.ID()] = product{
			Name:  p.Name(),
			Price: p.Price(),
		}
	}
	return nil
}

func (r *InMemoryRepo) GetProducts(_ context.Context, ids []uuid.UUID) ([]products.Product, error) {
	ps := make([]products.Product, 0, len(ids))
	for _, id := range ids {
		p, ok := r.products[id]
		if !ok {
			return nil, fmt.Errorf("%w: id %s", products.ErrProductNotFound, id)
		}
		product, err := products.NewProduct(id, p.Name, p.Price)
		if err != nil {
			return nil, err
		}
		ps = append(ps, *product)
	}
	return ps, nil
}
