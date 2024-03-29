package orders

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var (
	ErrInvalidItem = errors.New("invalid order item")
)

type Item struct {
	id       uuid.UUID
	name     string
	price    float64
	quantity int
}

func NewItem(id uuid.UUID, name string, price float64, quantity int) (*Item, error) {
	if quantity <= 0 {
		return nil, fmt.Errorf("%w: invalid quantity: %d", ErrInvalidItem, quantity)
	}
	if name == "" {
		return nil, fmt.Errorf("%w: invalid name", ErrInvalidItem)
	}
	if price <= 0 {
		return nil, fmt.Errorf("%w: invalid price: %f", ErrInvalidItem, price)
	}

	return &Item{
		id:       id,
		name:     name,
		price:    price,
		quantity: quantity,
	}, nil
}

func (i *Item) ID() uuid.UUID {
	return i.id
}

func (i *Item) Name() string {
	return i.name
}

func (i *Item) Price() float64 {
	return i.price
}

func (i *Item) Quantity() int {
	return i.quantity
}
