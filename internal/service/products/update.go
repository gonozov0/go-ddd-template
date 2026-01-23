package products

import (
	"context"

	"fmt"

	"go-ddd-template/internal/domain/products"
	"go-ddd-template/internal/domain/shared/valueobjects"
)

func (s ProductService) GetImageUploadURL(
	ctx context.Context,
	id valueobjects.ProductID,
	imageFilename valueobjects.ImageFilename,
) (imageUploadURL valueobjects.ImageURL, err error) {
	if err = s.repo.UpdateProduct(ctx, id, func(product *products.Product) error {
		if !product.GetImageFilename().IsEmpty() && !product.GetImageURL().IsEmpty() {
			if err := s.imageStorage.DeleteImage(ctx, product.GetImageFilename()); err != nil {
				return fmt.Errorf("failed to delete old image %s: %w", product.GetImageFilename().String(), err)
			}
		}

		product.SetImageFilename(imageFilename)

		imageUploadURL, err = s.imageStorage.GetImageUploadURL(imageFilename)
		if err != nil {
			return fmt.Errorf("failed to get image upload url: %w", err)
		}

		return nil
	}); err != nil {
		return valueobjects.EmptyImageURL, fmt.Errorf("failed to update product: %w", err)
	}

	return imageUploadURL, nil
}

func (s ProductService) ConfirmImageUpload(
	ctx context.Context,
	id valueobjects.ProductID,
) (imagePublicURL valueobjects.ImageURL, err error) {
	if err = s.repo.UpdateProduct(ctx, id, func(product *products.Product) error {
		imagePublicURL, err = s.imageStorage.ConfirmImageUpload(ctx, product.GetImageFilename())
		if err != nil {
			return fmt.Errorf("failed to confirm image upload: %w", err)
		}

		product.SetImageURL(imagePublicURL)

		return nil
	}); err != nil {
		return valueobjects.EmptyImageURL, fmt.Errorf("failed to update product: %w", err)
	}

	return imagePublicURL, nil
}
