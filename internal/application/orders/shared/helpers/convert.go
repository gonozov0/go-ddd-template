package helpers

import (
	pb "go-ddd-template/generated/server"
	domain "go-ddd-template/internal/domain/orders"
)

func ToCreateOrderRequest(order domain.Order) *pb.CreateOrderRequest {
	return &pb.CreateOrderRequest{
		ProductIds: order.GetProductIDs().Strings(),
	}
}
