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
	"go-ddd-template/pkg/auth"
	"go-ddd-template/pkg/grpcutils"
	loggerutils "go-ddd-template/pkg/logger/utils"
)

func (h OrderHandlers) DeleteOrder(
	ctx context.Context,
	req *pb.DeleteOrderRequest,
) (*emptypb.Empty, error) {
	userInfo := auth.GetUserInfo(ctx)
	if userInfo.IsEmpty() {
		return nil, auth.ErrMissingUserMetadata
	}

	if !userInfo.IsAdmin() {
		return nil, auth.ErrNoAdminRole
	}

	orderID, err := valueobjects.NewOrderIDFromString(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = h.orderService.DeleteOrder(ctx, orderID)
	if err != nil {
		slog.ErrorContext(
			ctx,
			"failed to delete order",
			loggerutils.ErrAttr(fmt.Errorf("failed to delete order: %w", err)),
		)

		return nil, grpcutils.ErrInternalError
	}

	return &emptypb.Empty{}, nil
}
