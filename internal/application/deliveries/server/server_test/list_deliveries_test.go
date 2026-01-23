package server_test

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "go-ddd-template/generated/server"
	deliveryunithelpers "go-ddd-template/internal/application/deliveries/shared/helpers"
	domain "go-ddd-template/internal/domain/deliveries"
	"go-ddd-template/pkg/grpcutils"
)

func (s *DeliveriesSuite) TestListDeliveries() {
	s.Run("successful list empty deliveries", func() {
		resp, err := s.GRPCHandlers.ListDeliveries(context.Background(), &emptypb.Empty{})
		grpcutils.CheckCodeWithSuite(s, err, codes.OK)

		s.Require().NotNil(resp)
		s.Require().Empty(resp.Deliveries)
	})
	s.Run("successful list with deliveries", func() {
		deliveries := []domain.Delivery{
			*deliveryunithelpers.CreateDelivery(s, s.DeliveriesRepo),
			*deliveryunithelpers.CreateDelivery(s, s.DeliveriesRepo),
			*deliveryunithelpers.CreateDelivery(s, s.DeliveriesRepo),
		}

		resp, err := s.GRPCHandlers.ListDeliveries(context.Background(), &emptypb.Empty{})
		grpcutils.CheckCodeWithSuite(s, err, codes.OK)

		expectedDeliveries := make([]*pb.Delivery, 0, len(deliveries))
		for _, delivery := range deliveries {
			expectedDeliveries = append(expectedDeliveries, &pb.Delivery{
				Id:        delivery.GetID().String(),
				OrderId:   delivery.GetOrderID().String(),
				CreatedAt: delivery.GetCreatedAt().Format(time.RFC3339),
			})
		}

		s.Require().ElementsMatch(expectedDeliveries, resp.Deliveries)
	})
}
