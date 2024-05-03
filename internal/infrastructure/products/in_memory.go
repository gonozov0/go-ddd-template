package products

import (
	"fmt"
	"sync"

	"go-echo-ddd-template/internal/domain/products"

	"github.com/google/uuid"
)

type product struct {
	Name      string
	Price     float64
	Available bool
}

type InMemoryRepo struct {
	products map[uuid.UUID]product
	tx       sync.Mutex
}

func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{
		products: make(map[uuid.UUID]product),
	}
}

func (r *InMemoryRepo) SaveProducts(ps []products.Product) error {
	defer r.tx.Unlock()

	for _, p := range ps {
		r.products[p.ID()] = product{
			Name:      p.Name(),
			Price:     p.Price(),
			Available: p.Available(),
		}
	}
	return nil
}

func (r *InMemoryRepo) GetProductsForUpdate(ids []uuid.UUID) ([]products.Product, error) {
	r.tx.Lock()

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

func (r *InMemoryRepo) CancelUpdate() {
	r.tx.TryLock()
	r.tx.Unlock()
}
