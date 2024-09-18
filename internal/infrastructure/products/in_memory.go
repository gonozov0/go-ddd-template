package products

import (
	"context"
	"fmt"
	"sync"

	"go-echo-template/internal/domain/products"

	"github.com/google/uuid"
)

type product struct {
	Name      string
	Price     float64
	Available bool
}

type InMemoryRepo struct {
	products map[uuid.UUID]product
	mu       sync.Mutex
}

func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{
		products: make(map[uuid.UUID]product),
	}
}

func (r *InMemoryRepo) SaveProducts(_ context.Context, ps []products.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, p := range ps {
		r.products[p.ID()] = product{
			Name:      p.Name(),
			Price:     p.Price(),
			Available: p.Available(),
		}
	}
	return nil
}

func (r *InMemoryRepo) GetProductsForUpdate(_ context.Context, ids []uuid.UUID) ([]products.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	ps := make([]products.Product, 0, len(ids))
	for _, id := range ids {
		p, ok := r.products[id]
		if !ok {
			return nil, fmt.Errorf("%w: id %s", products.ErrProductNotFound, id)
		}
		product, err := products.NewProduct(id, p.Name, p.Price, p.Available)
		if err != nil {
			return nil, err
		}
		ps = append(ps, *product)
	}
	return ps, nil
}
