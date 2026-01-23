package helpers

import (
	"context"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "go-ddd-template/generated/server"
	"go-ddd-template/integration_tests/suites"
	"go-ddd-template/internal/domain/shared/valueobjects"
	"go-ddd-template/pkg/grpcutils"
	loggerutils "go-ddd-template/pkg/logger/utils"
)

type HelperSuite struct {
	baseSuite *suites.BaseSuite

	grpcClient pb.OrderServiceClient
}

func NewHelperSuite(baseSuite *suites.BaseSuite) *HelperSuite {
	handlerSuite := &HelperSuite{
		baseSuite: baseSuite,

		grpcClient: pb.NewOrderServiceClient(baseSuite.ServerConn),
	}

	return handlerSuite
}

func (s *HelperSuite) CreateOrder(
	userCtx context.Context,
	productIDs valueobjects.ProductIDs,
) valueobjects.OrderID {
	req := &pb.CreateOrderRequest{
		ProductIds: productIDs.Strings(),
	}

	resp, err := s.grpcClient.CreateOrder(userCtx, req)
	grpcutils.CheckCodeWithSuite(s.baseSuite, err, codes.OK)
	s.baseSuite.Require().NotNil(resp)

	orderID, err := valueobjects.NewOrderIDFromString(resp.GetId())
	s.baseSuite.Require().NoError(err)

	s.baseSuite.T().Cleanup(func() { s.DeleteOrderInCleanup(orderID) })

	return orderID
}

func (s *HelperSuite) DeleteOrderInCleanup(
	orderID valueobjects.OrderID,
) {
	if orderID == valueobjects.EmptyOrderID {
		return
	}

	_, err := s.grpcClient.DeleteOrder(s.baseSuite.UserCtx, &pb.DeleteOrderRequest{
		Id: orderID.String(),
	})

	grpcStatus, ok := status.FromError(err)
	if err != nil || !ok || grpcStatus.Code() != codes.OK {
		slog.Error(
			"failed to delete order",
			loggerutils.ErrAttr(err),
			loggerutils.Attr("grpc_status", grpcStatus.Code().String()),
			loggerutils.Attr("grpc_message", grpcStatus.Message()),
			loggerutils.Attr("order_id", orderID.String()),
		)

		return
	}

	slog.Info("Deleted order", loggerutils.Attr("order_id", orderID.String()))
}
