package application

import (
	"go-ddd-template/internal/service/deliveries"
	"go-ddd-template/internal/service/orders"
	"go-ddd-template/internal/service/products"
	"go-ddd-template/internal/service/users"
)

type UserRepository interface {
	users.UserRepository
}

type OrderRepository interface {
	orders.OrderRepository
}

type ProductRepository interface {
	products.ProductRepository
}

type DeliveryRepository interface {
	deliveries.DeliveryRepository
}
