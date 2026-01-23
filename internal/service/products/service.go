package products

import (
	"context"

	"go-ddd-template/internal/domain/products"
	"go-ddd-template/internal/domain/shared/valueobjects"
)

type ProductRepository interface {
	CreateProducts(ctx context.Context, createFn func() (products.Products, error)) error
	GetProducts(ctx context.Context, ids valueobjects.ProductIDs) ([]products.Product, error)
	DeleteProducts(ctx context.Context, ids valueobjects.ProductIDs, deleteFn func([]products.Product) error) error
	UpdateProduct(ctx context.Context, id valueobjects.ProductID, updateFn func(*products.Product) error) error
}

type ImageStorage interface {
	GetImageUploadURL(filename valueobjects.ImageFilename) (valueobjects.ImageURL, error)
	ConfirmImageUpload(ctx context.Context, filename valueobjects.ImageFilename) (valueobjects.ImageURL, error)
	DeleteImage(ctx context.Context, filename valueobjects.ImageFilename) error
}

type ProductService struct {
	repo         ProductRepository
	imageStorage ImageStorage
}

func NewProductService(r ProductRepository, imageStorage ImageStorage) ProductService {
	return ProductService{
		repo:         r,
		imageStorage: imageStorage,
	}
}
