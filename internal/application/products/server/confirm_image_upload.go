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
	imagestorage "go-ddd-template/pkg/image_storage"
	loggerutils "go-ddd-template/pkg/logger/utils"
)

func (h ProductHandlers) ConfirmProductImageUpload(
	ctx context.Context,
	req *pb.ConfirmProductImageUploadRequest,
) (*pb.ConfirmProductImageUploadResponse, error) {
	productID, err := valueobjects.NewProductIDFromString(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", err.Error())
	}

	imageURL, err := h.productService.ConfirmImageUpload(ctx, productID)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			return nil, status.Errorf(codes.NotFound, "%s", err.Error())
		}

		if errors.Is(err, imagestorage.ErrImageNotFound) {
			return nil, status.Errorf(codes.FailedPrecondition, "%s", err.Error())
		}

		slog.ErrorContext(
			ctx,
			"failed to confirm product image upload",
			loggerutils.ErrAttr(fmt.Errorf("failed to confirm product image upload: %w", err)),
		)

		return nil, grpcutils.ErrInternalError
	}

	return &pb.ConfirmProductImageUploadResponse{
		ImageUrl: imageURL.String(),
	}, nil
}
