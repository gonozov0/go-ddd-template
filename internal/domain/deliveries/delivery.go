package deliveries

import (
	"github.com/google/uuid"

	"go-ddd-template/internal/domain/shared/valueobjects"
)

type Delivery struct {
	id        valueobjects.DeliveryID
	orderID   valueobjects.OrderID
	createdAt valueobjects.Timestamp
}

func NewDelivery(
	deliveryID valueobjects.DeliveryID,
	orderID valueobjects.OrderID,
	createdAt valueobjects.Timestamp,
) *Delivery {
	return &Delivery{
		id:        deliveryID,
		orderID:   orderID,
		createdAt: createdAt,
	}
}

func CreateDelivery(orderID valueobjects.OrderID) *Delivery {
	return NewDelivery(
		valueobjects.DeliveryID(uuid.New()),
		orderID,
		valueobjects.NewTimestampNow(),
	)
}

func (d *Delivery) GetID() valueobjects.DeliveryID {
	return d.id
}

func (d *Delivery) GetOrderID() valueobjects.OrderID {
	return d.orderID
}

func (d *Delivery) GetCreatedAt() valueobjects.Timestamp {
	return d.createdAt
}
