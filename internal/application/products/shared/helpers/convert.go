package helpers

import (
	"fmt"

	pb "go-ddd-template/generated/server"
	domain "go-ddd-template/internal/domain/products"
	"go-ddd-template/internal/domain/shared/valueobjects"
)

func ToCreateProductsRequest(products domain.Products) *pb.CreateProductsRequest {
	var req pb.CreateProductsRequest

	for _, product := range products {
		req.Items = append(req.Items, &pb.ProductRequestItem{
			Name:  product.GetName().String(),
			Price: product.GetPrice().Float64(),
		})
	}

	return &req
}

func UpdateProductsWithIDs(productIDs valueobjects.ProductIDs, products domain.Products) (domain.Products, error) {
	if len(productIDs) != len(products) {
		return nil, fmt.Errorf("length of productIDs and products must be equal")
	}

	updatedProducts := make(domain.Products, 0, len(products))
	for i, product := range products {
		updatedProducts = append(updatedProducts, *domain.NewProduct(
			productIDs[i],
			product.GetName(),
			product.GetPrice(),
			product.GetStatus(),
			product.GetImageFilename(),
			product.GetImageURL(),
		))
	}

	return updatedProducts, nil
}
