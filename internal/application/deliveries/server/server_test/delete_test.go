package server_test

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"

	pb "go-ddd-template/generated/server"
	deliveryunithelpers "go-ddd-template/internal/application/deliveries/shared/helpers"
	"go-ddd-template/pkg/grpcutils"
)

func (s *DeliveriesSuite) TestDeleteDelivery() {
	s.Run("Failed due to invalid delivery uuid", func() {
		_, err := s.GRPCHandlers.DeleteDelivery(context.Background(), &pb.DeleteDeliveryRequest{
			Id: "inalid-uuid",
		})
		grpcutils.CheckCodeWithSuite(s, err, codes.InvalidArgument)

		s.Require().ErrorContains(err, "invalid delivery id")
	})
	s.Run("Successfully deleted existing delivery", func() {
		delivery := deliveryunithelpers.CreateDelivery(s, s.DeliveriesRepo)

		_, err := s.GRPCHandlers.DeleteDelivery(context.Background(), &pb.DeleteDeliveryRequest{
			Id: delivery.GetID().String(),
		})
		grpcutils.CheckCodeWithSuite(s, err, codes.OK)

		deliveries, err := s.DeliveriesRepo.ListDeliveries(context.Background())
		s.Require().NoError(err)
		s.Require().Len(deliveries, 0, "delivery should be deleted but still exists")
	})
	s.Run("Successfully deleted non-existing delivery", func() {
		randomID := uuid.New()

		_, err := s.GRPCHandlers.DeleteDelivery(context.Background(), &pb.DeleteDeliveryRequest{
			Id: randomID.String(),
		})
		grpcutils.CheckCodeWithSuite(s, err, codes.OK)
	})
}
