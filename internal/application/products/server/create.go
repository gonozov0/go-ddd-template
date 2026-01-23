package products

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "go-ddd-template/generated/server"
	domain "go-ddd-template/internal/domain/products"
	"go-ddd-template/internal/domain/shared/valueobjects"
	service "go-ddd-template/internal/service/products"
	"go-ddd-template/pkg/grpcutils"
	loggerutils "go-ddd-template/pkg/logger/utils"
)

func (h ProductHandlers) CreateProducts(
	ctx context.Context,
	req *pb.CreateProductsRequest,
) (*pb.CreateProductsResponse, error) {
	productsToCreate := make(service.ProductsToCreate, 0, len(req.GetItems()))

	for _, item := range req.GetItems() {
		name, err := valueobjects.NewProductName(item.GetName())
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "%s", err.Error())
		}

		price, err := valueobjects.NewProductPrice(float64(item.GetPrice()))
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "%s", err.Error())
		}

		productsToCreate = append(productsToCreate, service.ProductToCreate{
			Name:  name,
			Price: price,
		})
	}

	productIds, err := h.productService.CreateProducts(ctx, productsToCreate)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidProduct) {
			return nil, status.Errorf(codes.InvalidArgument, "%s", err.Error())
		}

		slog.ErrorContext(
			ctx,
			"failed to create products",
			loggerutils.ErrAttr(fmt.Errorf("failed to create products: %w", err)),
		)

		return nil, grpcutils.ErrInternalError
	}

	return &pb.CreateProductsResponse{
		Ids: productIds.Strings(),
	}, nil
}
