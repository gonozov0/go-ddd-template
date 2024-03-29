package products

import (
	"sync"

	"go-echo-ddd-template/internal/domain/products"

	"github.com/google/uuid"
)

type product struct {
	Name     string
	Price    float64
	Category string
	Quantity int
}

type InMemoryRepo struct {
	products map[uuid.UUID]product
	mu       sync.RWMutex
}

func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{
		products: make(map[uuid.UUID]product),
	}
}

func (r *InMemoryRepo) SaveProduct(p products.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// TODO: check mutex for each product to simulate transaction with blocking raws
	r.products[p.ID()] = product{
		Name:     p.Name(),
		Price:    p.Price(),
		Category: p.Category(),
		Quantity: p.Quantity(),
	}
	return nil
}

func (r *InMemoryRepo) SaveProducts(ps []products.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// TODO: check mutex for each product to simulate transaction with blocking raws
	for _, p := range ps {
		r.products[p.ID()] = product{
			Name:     p.Name(),
			Price:    p.Price(),
			Category: p.Category(),
			Quantity: p.Quantity(),
		}
	}
	return nil
}

func (r *InMemoryRepo) GetProduct(id uuid.UUID) (*products.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, ok := r.products[id]
	if !ok {
		return nil, products.ErrProductNotFound
	}
	product, err := products.NewProduct(id, p.Name, p.Price, p.Category, p.Quantity)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (r *InMemoryRepo) GetProductsForUpdate(ids []uuid.UUID) ([]products.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// TODO: make mutex for each product to simulate transaction with blocking raws
	ps := make([]products.Product, 0, len(ids))
	for _, id := range ids {
		p, ok := r.products[id]
		if !ok {
			return nil, products.ErrProductNotFound
		}
		product, err := products.NewProduct(id, p.Name, p.Price, p.Category, p.Quantity)
		if err != nil {
			return nil, err
		}
		ps = append(ps, *product)
	}
	return ps, nil
}
