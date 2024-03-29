package orders

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrOrderNotFound = errors.New("order not found")
)

type Order struct {
	id     uuid.UUID
	userID uuid.UUID
	status OrderStatus
	items  []Item
}

func NewOrder(id uuid.UUID, userID uuid.UUID, status OrderStatus, items []Item) (*Order, error) {
	return &Order{
		id:     id,
		userID: userID,
		status: status,
		items:  items,
	}, nil
}

func CreateOrder(userID uuid.UUID, items []Item) (*Order, error) {
	return NewOrder(uuid.New(), userID, OrderStatusCreated, items)
}

func (o *Order) ID() uuid.UUID {
	return o.id
}

func (o *Order) UserID() uuid.UUID {
	return o.userID
}

func (o *Order) Status() OrderStatus {
	return o.status
}

func (o *Order) Items() []Item {
	return o.items
}

func (o *Order) Price() float64 {
	var total float64
	for _, item := range o.items {
		total += item.Price() * float64(item.Quantity())
	}
	return total
}

func (o *Order) Pay() error {
	return errors.New("not implemented")
}

func (o *Order) MakeInvoice() (string, error) {
	if o.status != OrderStatusPaid {
		return "", errors.New("order is not paid")
	}
	return "", errors.New("not implemented")
}
