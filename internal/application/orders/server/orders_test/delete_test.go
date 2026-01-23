package orders_test

import (
	"context"

	"google.golang.org/grpc/codes"

	pb "go-ddd-template/generated/server"
	domain "go-ddd-template/internal/domain/orders"
	"go-ddd-template/pkg/auth"
	"go-ddd-template/pkg/grpcutils"
)

func (s *OrdersSuite) TestDeleteOrder() {
	order := s.createOrder()

	s.Run("Failed due to unauthenticated request", func() {
		_, err := s.GRPCHandlers.DeleteOrder(context.Background(), &pb.DeleteOrderRequest{
			Id: order.GetID().String(),
		})
		grpcutils.CheckCodeWithSuite(s, err, codes.Unauthenticated)

		s.Require().ErrorContains(err, auth.ErrMissingUserMetadata.Error())
	})

	s.Run("Failed due to invalid order uuid", func() {
		_, err := s.GRPCHandlers.DeleteOrder(s.AdminCtx, &pb.DeleteOrderRequest{
			Id: "invalid-uuid",
		})
		grpcutils.CheckCodeWithSuite(s, err, codes.InvalidArgument)

		s.Require().ErrorContains(err, "invalid order id")
	})

	s.Run("Successfully deleted with valid state", func() {
		_, err := s.GRPCHandlers.DeleteOrder(s.AdminCtx, &pb.DeleteOrderRequest{
			Id: order.GetID().String(),
		})
		grpcutils.CheckCodeWithSuite(s, err, codes.OK)

		_, err = s.OrdersRepo.GetOrder(context.Background(), order.GetID())
		s.Require().ErrorIs(err, domain.ErrOrderNotFound)
	})
}
