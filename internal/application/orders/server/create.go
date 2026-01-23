package server

import (
	"context"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"errors"
	"fmt"

	pb "go-ddd-template/generated/server"
	ordersdomain "go-ddd-template/internal/domain/orders"
	"go-ddd-template/internal/domain/shared/valueobjects"
	usersdomain "go-ddd-template/internal/domain/users"
	ordersservice "go-ddd-template/internal/service/orders"
	"go-ddd-template/pkg/auth"
	"go-ddd-template/pkg/grpcutils"
	loggerutils "go-ddd-template/pkg/logger/utils"
)

func (h OrderHandlers) CreateOrder(
	ctx context.Context,
	req *pb.CreateOrderRequest,
) (*pb.CreateOrderResponse, error) {
	userInfo := auth.GetUserInfo(ctx)
	if userInfo.IsEmpty() {
		return nil, auth.ErrMissingUserMetadata
	}

	userID, err := valueobjects.NewUserID(userInfo.ID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user id")
	}

	productIDs, err := valueobjects.NewProductIDsFromStrings(req.GetProductIds())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid product ids")
	}

	order, err := h.orderService.CreateOrder(ctx, ordersservice.OrderToCreate{
		UserID:     userID,
		ProductIDs: productIDs,
	})
	if err != nil {
		return nil, parseCreateOrderError(ctx, err)
	}

	return &pb.CreateOrderResponse{
		Id: order.GetID().String(),
	}, nil
}

func parseCreateOrderError(ctx context.Context, err error) error {
	switch {
	case errors.Is(err, ordersdomain.ErrProductAlreadyPublished):
		return status.Errorf(codes.Aborted, "%s", ordersdomain.ErrProductAlreadyPublished.Error())
	case errors.Is(err, usersdomain.ErrUserNotFound):
		return status.Errorf(codes.NotFound, "%s", usersdomain.ErrUserNotFound.Error())
	case errors.Is(err, ordersdomain.ErrProductNotFound):
		return status.Errorf(codes.NotFound, "%s", ordersdomain.ErrProductNotFound.Error())
	case errors.Is(err, ordersdomain.ErrReserveInitedProduct):
		return status.Errorf(codes.NotFound, "%s", ordersdomain.ErrReserveInitedProduct.Error())
	default:
		slog.ErrorContext(
			ctx,
			"failed to create order",
			loggerutils.ErrAttr(fmt.Errorf("failed to create order: %w", err)),
		)

		return grpcutils.ErrInternalError
	}
}
