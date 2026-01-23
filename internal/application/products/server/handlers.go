package products

import (
	protobuf "go-ddd-template/generated/server"
	service "go-ddd-template/internal/service/products"
)

type ProductHandlers struct {
	protobuf.UnimplementedProductServiceServer
	productService service.ProductService
	imageStorage   service.ImageStorage
}

func SetupHandlers(productRepo service.ProductRepository, imageStorage service.ImageStorage) ProductHandlers {
	productService := service.NewProductService(productRepo, imageStorage)

	//nolint:exhaustivestruct
	return ProductHandlers{
		productService: productService,
		imageStorage:   imageStorage,
	}
}
