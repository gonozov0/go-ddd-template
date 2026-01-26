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

func NewOrderStatus(status string) (OrderStatus, error) {
	switch OrderStatus(status) {
	case OrderStatusCreated:
		return OrderStatusCreated, nil
	case OrderStatusProcessing:
		return OrderStatusProcessing, nil
	default:
		return "", fmt.Errorf("invalid order status: %s", status)
	}
}

func (s OrderStatus) String() string {
	return string(s)
}
