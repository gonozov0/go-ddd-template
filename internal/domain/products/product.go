package products

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var (
	ErrProductNotFound            = errors.New("product not found")
	ErrInvalidQuantity            = errors.New("invalid quantity")
	ErrInvalidProduct             = errors.New("invalid product")
	ErrProductQuantityIsNotEnough = errors.New("product quantity is not enough")
)

type Product struct {
	id       uuid.UUID
	name     string
	price    float64
	category string // TODO: implement category sub aggregate
	quantity int
}

func NewProduct(id uuid.UUID, name string, price float64, category string, quantity int) (*Product, error) {
	if err := validateProductName(name); err != nil {
		return nil, err
	}
	if err := validateProductPrice(price); err != nil {
		return nil, err
	}
	if err := validateProductCategory(category); err != nil {
		return nil, err
	}
	if err := validateProductQuantity(quantity); err != nil {
		return nil, err
	}

	return &Product{
		id:       id,
		name:     name,
		price:    price,
		category: category,
		quantity: quantity,
	}, nil
}

func CreateProduct(name string, price float64, category string) (*Product, error) {
	return NewProduct(uuid.New(), name, price, category, 0)
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

func (p *Product) Category() string {
	return p.category
}

func (p *Product) Quantity() int {
	return p.quantity
}

func (p *Product) AddQuantity(quantity int) error {
	if quantity <= 0 {
		return fmt.Errorf("%w: quantity must be greater than 0", ErrInvalidQuantity)
	}
	p.quantity += quantity
	return nil
}

func (p *Product) ReduceQuantity(quantity int) error {
	if p.quantity < quantity {
		return fmt.Errorf(
			"%w: product id %s, need %d, but only %d available",
			ErrProductQuantityIsNotEnough,
			p.id,
			quantity,
			p.quantity,
		)
	}
	p.quantity -= quantity
	return nil
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

func validateProductCategory(category string) error {
	if category == "" {
		return fmt.Errorf("%w: category is required", ErrInvalidProduct)
	}
	return nil
}

func validateProductQuantity(quantity int) error {
	if quantity < 0 {
		return fmt.Errorf("%w: quantity must be greater than 0", ErrInvalidProduct)
	}
	return nil
}
