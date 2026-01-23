package server

import (
	"context"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"fmt"

	pb "go-ddd-template/generated/server"
	"go-ddd-template/internal/domain/shared/valueobjects"
	"go-ddd-template/pkg/grpcutils"
	loggerutils "go-ddd-template/pkg/logger/utils"
)

func (h DeliveryHandlers) DeleteDelivery(
	ctx context.Context,
	req *pb.DeleteDeliveryRequest,
) (*emptypb.Empty, error) {
	deliveryID, err := valueobjects.NewDeliveryIDFromString(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, valueobjects.ErrInvalidDeliveryID.Error())
	}

	err = h.deliveryService.DeleteDelivery(ctx, deliveryID)
	if err != nil {
		slog.ErrorContext(
			ctx,
			"failed to delete delivery",
			loggerutils.ErrAttr(fmt.Errorf("failed to delete delivery: %w", err)),
		)

		return nil, grpcutils.ErrInternalError
	}

	return &emptypb.Empty{}, nil
}
