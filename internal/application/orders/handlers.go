package orders

import (
	"go-echo-template/generated/protobuf"
	service "go-echo-template/internal/service/orders"
)

type OrderHandlers struct {
	protobuf.UnimplementedOrderServiceServer
	orderRepo   service.OrderRepository
	userRepo    service.UserRepository
	productRepo service.ProductRepository
}

func SetupHandlers(or service.OrderRepository, ur service.UserRepository, pr service.ProductRepository) OrderHandlers {
	return OrderHandlers{
		orderRepo:   or,
		userRepo:    ur,
		productRepo: pr,
	}
}
