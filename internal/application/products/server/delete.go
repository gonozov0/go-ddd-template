package products

import (
	"context"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"fmt"

	protobuf "go-ddd-template/generated/server"
	"go-ddd-template/internal/domain/shared/valueobjects"
	"go-ddd-template/pkg/grpcutils"
	loggerutils "go-ddd-template/pkg/logger/utils"
)

func (h ProductHandlers) DeleteProducts(
	ctx context.Context,
	req *protobuf.DeleteProductsRequest,
) (*emptypb.Empty, error) {
	productIDs, err := valueobjects.NewProductIDsFromStrings(req.GetIds())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := h.productService.DeleteProducts(ctx, productIDs); err != nil {
		slog.ErrorContext(
			ctx,
			"failed to delete products",
			loggerutils.ErrAttr(fmt.Errorf("failed to delete products: %w", err)),
		)

		return nil, grpcutils.ErrInternalError
	}

	return &emptypb.Empty{}, nil
}
