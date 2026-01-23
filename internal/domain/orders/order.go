package orders

import (
	"github.com/google/uuid"

	"errors"

	"go-ddd-template/internal/domain/shared/valueobjects"
)

var (
	ErrOrderNotFound          = errors.New("order not found")
	ErrOrderAlreadyProcessing = errors.New("order already processing")
)

type Order struct {
	id       valueobjects.OrderID
	userID   valueobjects.UserID
	status   OrderStatus
	products []Product
}

func NewOrder(
	id valueobjects.OrderID,
	userID valueobjects.UserID,
	status OrderStatus,
	products []Product,
) *Order {
	return &Order{
		id:       id,
		userID:   userID,
		status:   status,
		products: products,
	}
}

func CreateOrder(userID valueobjects.UserID, products []Product) *Order {
	return NewOrder(
		valueobjects.OrderID(uuid.New()),
		userID,
		OrderStatusCreated,
		products,
	)
}

func (o *Order) GetID() valueobjects.OrderID {
	return o.id
}

func (o *Order) GetUserID() valueobjects.UserID {
	return o.userID
}

func (o *Order) GetStatus() OrderStatus {
	return o.status
}

func (o *Order) GetProductIDs() valueobjects.ProductIDs {
	ids := make(valueobjects.ProductIDs, 0, len(o.products))

	for _, product := range o.products {
		ids = append(ids, product.GetID())
	}

	return ids
}

func (o *Order) GetPrice() valueobjects.ProductPrice {
	var total valueobjects.ProductPrice

	for _, product := range o.products {
		total += product.GetPrice()
	}

	return total
}

func (o *Order) Process() error {
	if o.status == OrderStatusProcessing {
		return ErrOrderAlreadyProcessing
	}

	o.status = OrderStatusProcessing

	return nil
}
