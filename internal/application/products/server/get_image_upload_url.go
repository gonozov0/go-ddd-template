package products

import (
	"context"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"errors"
	"fmt"

	pb "go-ddd-template/generated/server"
	domain "go-ddd-template/internal/domain/products"
	"go-ddd-template/internal/domain/shared/valueobjects"
	"go-ddd-template/pkg/grpcutils"
	loggerutils "go-ddd-template/pkg/logger/utils"
)

func (h ProductHandlers) GetProductImageUploadURL(
	ctx context.Context,
	req *pb.GetProductImageUploadURLRequest,
) (*pb.GetProductImageUploadURLResponse, error) {
	productID, err := valueobjects.NewProductIDFromString(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", err.Error())
	}

	imageFilename, err := valueobjects.NewRequiredImageFilenameForProduct(req.GetFilename(), productID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", err.Error())
	}

	presignedURL, err := h.productService.GetImageUploadURL(ctx, productID, imageFilename)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			return nil, status.Errorf(codes.NotFound, "%s", err.Error())
		}

		slog.ErrorContext(
			ctx,
			"failed to get product image upload URL",
			loggerutils.ErrAttr(fmt.Errorf("failed to get product image upload URL: %w", err)),
		)

		return nil, grpcutils.ErrInternalError
	}

	return &pb.GetProductImageUploadURLResponse{
		UploadUrl: presignedURL.String(),
	}, nil
}
