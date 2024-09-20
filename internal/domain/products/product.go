package products

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var (
	ErrProductNotFound = errors.New("product not found")
	ErrInvalidProduct  = errors.New("invalid product")
)

type Product struct {
	id    uuid.UUID
	name  string
	price float64
}

func NewProduct(id uuid.UUID, name string, price float64) (*Product, error) {
	if err := validateProductName(name); err != nil {
		return nil, err
	}
	if err := validateProductPrice(price); err != nil {
		return nil, err
	}

	return &Product{
		id:    id,
		name:  name,
		price: price,
	}, nil
}

func CreateProduct(name string, price float64) (*Product, error) {
	return NewProduct(uuid.New(), name, price)
}

func (p *Product) ID() uuid.UUID {
	return p.id
}

func (p *Product) Name() string {
	return p.name
}

func (p *Product) Price() float64 {
	return p.price
}

func validateProductName(name string) error {
	if name == "" {
		return fmt.Errorf("%w: name is required", ErrInvalidProduct)
	}
	return nil
}

func validateProductPrice(price float64) error {
	if price <= 0 {
		return fmt.Errorf("%w: price must be greater than 0", ErrInvalidProduct)
	}
	return nil
}
