package helpers

import (
	"context"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"fmt"

	pb "go-ddd-template/generated/server"
	"go-ddd-template/integration_tests/suites"
	deliveriesdomain "go-ddd-template/internal/domain/deliveries"
	"go-ddd-template/internal/domain/shared/valueobjects"
	"go-ddd-template/pkg/backoff"
	loggerutils "go-ddd-template/pkg/logger/utils"
)

type HelperSuite struct {
	baseSuite *suites.BaseSuite

	grpcClient pb.DeliveryServiceClient
}

func NewHelperSuite(baseSuite *suites.BaseSuite) *HelperSuite {
	handlerSuite := &HelperSuite{
		baseSuite: baseSuite,

		grpcClient: pb.NewDeliveryServiceClient(baseSuite.ServerConn),
	}

	return handlerSuite
}

func (s *HelperSuite) WaitForDeliveryCreation(
	orderID valueobjects.OrderID,
) valueobjects.DeliveryID {
	var deliveryID valueobjects.DeliveryID

	err := backoff.RunWithRetry(context.Background(), func() error {
		listDeliveriesResp, err := s.grpcClient.ListDeliveries(s.baseSuite.UserCtx, &emptypb.Empty{})
		if err != nil {
			return err
		}

		for _, d := range listDeliveriesResp.Deliveries {
			if d.GetOrderId() == orderID.String() {
				deliveryID, err = valueobjects.NewDeliveryIDFromString(d.Id)
				if err != nil {
					return fmt.Errorf("failed to parse delivery id %s: %w", d.Id, err)
				}

				return nil
			}
		}

		return deliveriesdomain.NewRetriableDeliveryNotFoundError()
	}, 35)

	s.baseSuite.Require().NoError(err, "Delivery should be created for order %s", orderID.String())

	s.baseSuite.T().Cleanup(func() { s.DeleteDeliveryInCleanup(deliveryID) })

	return deliveryID
}

func (s *HelperSuite) DeleteDeliveryInCleanup(
	deliveryID valueobjects.DeliveryID,
) {
	if deliveryID.IsEmpty() {
		return
	}

	_, err := s.grpcClient.DeleteDelivery(s.baseSuite.AdminCtx, &pb.DeleteDeliveryRequest{
		Id: deliveryID.String(),
	})

	grpcStatus, ok := status.FromError(err)
	if err != nil || !ok || grpcStatus.Code() != codes.OK {
		slog.Error(
			"failed to delete delivery",
			loggerutils.ErrAttr(err),
			loggerutils.Attr("grpc_status", grpcStatus.Code().String()),
			loggerutils.Attr("grpc_message", grpcStatus.Message()),
			loggerutils.Attr("delivery_id", deliveryID.String()),
		)

		return
	}

	slog.Info("Deleted delivery", loggerutils.Attr("delivery_id", deliveryID.String()))
}
