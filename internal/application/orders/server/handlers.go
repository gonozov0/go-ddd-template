package server

import (
	pb "go-ddd-template/generated/server"
	service "go-ddd-template/internal/service/orders"
	"go-ddd-template/internal/service/users"
)

type OrderHandlers struct {
	pb.UnimplementedOrderServiceServer
	orderService service.OrderService
	userService  users.UserService
}

func SetupHandlers(or service.OrderRepository, ur users.UserRepository) OrderHandlers {
	userService := users.NewUserService(ur)
	orderService := service.NewOrderService(or, userService)

	//nolint:exhaustivestruct
	return OrderHandlers{
		orderService: orderService,
		userService:  userService,
	}
}
