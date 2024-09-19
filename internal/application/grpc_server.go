package application

import (
	"go-echo-template/generated/protobuf"
	"go-echo-template/internal/application/orders"
	"go-echo-template/internal/application/users"

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
