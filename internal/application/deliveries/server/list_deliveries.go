package server

import (
	"context"
	"time"

	"google.golang.org/protobuf/types/known/emptypb"

	"fmt"

	pb "go-ddd-template/generated/server"
)

func (h DeliveryHandlers) ListDeliveries(
	ctx context.Context,
	_ *emptypb.Empty,
) (*pb.ListDeliveriesResponse, error) {
	list, err := h.deliveryService.ListDeliveries(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list deliveries: %w", err)
	}

	//nolint:exhaustivestruct
	resp := &pb.ListDeliveriesResponse{}
	resp.Deliveries = make([]*pb.Delivery, 0, len(list))

	for _, d := range list {
		resp.Deliveries = append(resp.Deliveries, &pb.Delivery{
			Id:        d.GetID().String(),
			OrderId:   d.GetOrderID().String(),
			CreatedAt: d.GetCreatedAt().Format(time.RFC3339),
		})
	}

	return resp, nil
}
