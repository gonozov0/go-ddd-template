package orders

import (
	"fmt"
)

type OrderStatus string

const (
	OrderStatusCreated    OrderStatus = "created"
	OrderStatusProcessing OrderStatus = "processing"
)

var AllOrderStatuses = []OrderStatus{
	OrderStatusCreated,
	OrderStatusProcessing,
}

func (s OrderStatus) String() string {
	return string(s)
}

func NewOrderStatus(status string) (OrderStatus, error) {
	switch status {
	case OrderStatusCreated.String():
		return OrderStatusCreated, nil
	case OrderStatusProcessing.String():
		return OrderStatusProcessing, nil
	default:
		return "", fmt.Errorf("invalid order status: %s", status)
	}
}
