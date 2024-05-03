package products

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var (
	ErrProductNotFound        = errors.New("product not found")
	ErrInvalidProduct         = errors.New("invalid product")
	ErrProductAlreadyReserved = errors.New("product already reserved")
)

type Product struct {
	id        uuid.UUID
	name      string
	price     float64
	available bool
}

func NewProduct(id uuid.UUID, name string, price float64, available bool) (*Product, error) {
	if err := validateProductName(name); err != nil {
		return nil, err
	}
	if err := validateProductPrice(price); err != nil {
		return nil, err
	}

	return &Product{
		id:        id,
		name:      name,
		price:     price,
		available: available,
	}, nil
}

func CreateProduct(name string, price float64) (*Product, error) {
	return NewProduct(uuid.New(), name, price, true)
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

func (p *Product) Available() bool {
	return p.available
}

func (p *Product) Reserve() error {
	if !p.available {
		return ErrProductAlreadyReserved
	}
	p.available = false
	return nil
}

func (p *Product) CancelReservation() {
	p.available = true
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
