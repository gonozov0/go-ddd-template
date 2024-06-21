package application

import (
	"go-echo-ddd-template/generated/protobuf"
	"go-echo-ddd-template/internal/application/orders"
	"go-echo-ddd-template/internal/application/users"

	"google.golang.org/grpc"
)

type gRPCServer struct {
	users.UserHandlers
	orders.OrderHandlers
}

func SetupGRPCServer(userRepo UserRepository, orderRepo OrderRepository, productRepo ProductRepository) *grpc.Server {
	s := grpc.NewServer()

	server := gRPCServer{}
	server.UserHandlers = users.SetupHandlers(userRepo)
	server.OrderHandlers = orders.SetupHandlers(orderRepo, userRepo, productRepo)

	protobuf.RegisterOrderServiceServer(s, server)
	protobuf.RegisterUserServiceServer(s, server)

	return s
}
