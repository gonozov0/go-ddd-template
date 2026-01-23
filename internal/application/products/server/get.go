package products

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	protobuf "go-ddd-template/generated/server"
	domain "go-ddd-template/internal/domain/products"
	"go-ddd-template/internal/domain/shared/valueobjects"
	"go-ddd-template/pkg/grpcutils"
	loggerutils "go-ddd-template/pkg/logger/utils"
)

func (h ProductHandlers) GetProducts(
	ctx context.Context,
	req *protobuf.GetProductsRequest,
) (*protobuf.GetProductsResponse, error) {
	productIDs, err := valueobjects.NewProductIDsFromStrings(req.GetIds())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", err.Error())
	}

	products, err := h.productService.GetProducts(ctx, productIDs)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			return nil, status.Errorf(codes.NotFound, "%s", err.Error())
		}

		slog.ErrorContext(
			ctx,
			"failed to get domain",
			loggerutils.ErrAttr(fmt.Errorf("failed to get domain: %w", err)),
		)

		return nil, grpcutils.ErrInternalError
	}

	items := make([]*protobuf.ProductResponseItem, 0, len(products))
	for _, product := range products {
		items = append(items, &protobuf.ProductResponseItem{
			Id:     product.GetID().String(),
			Name:   product.GetName().String(),
			Price:  product.GetPrice().Float64(),
			Status: product.GetStatus().String(),
		})
	}

	return &protobuf.GetProductsResponse{
		Items: items,
	}, nil
}
