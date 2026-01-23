package products_test

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"

	pb "go-ddd-template/generated/server"
	productunithelpers "go-ddd-template/internal/application/products/shared/helpers"
	domain "go-ddd-template/internal/domain/products"
	"go-ddd-template/pkg/grpcutils"
)

func (s *ProductsSuite) TestGetProducts() {
	products := productunithelpers.CreateRandomProducts(s, s.ProductsRepo, 2)

	req := &pb.GetProductsRequest{Ids: products.IDs().Strings()}
	resp, err := s.GRPCHandlers.GetProducts(context.Background(), req)
	grpcutils.CheckCodeWithSuite(s, err, codes.OK)

	s.Require().Len(resp.GetItems(), 2)

	var expectedProducts []*pb.ProductResponseItem
	for _, product := range products {
		expectedProducts = append(expectedProducts, &pb.ProductResponseItem{
			Id:     product.GetID().String(),
			Name:   product.GetName().String(),
			Price:  product.GetPrice().Float64(),
			Status: product.GetStatus().String(),
		})
	}

	s.Require().ElementsMatch(expectedProducts, resp.Items)
}

func (s *ProductsSuite) TestGetProductsWithNonExistentId() {
	req := &pb.GetProductsRequest{Ids: []string{uuid.New().String()}}
	_, err := s.GRPCHandlers.GetProducts(context.Background(), req)
	message := grpcutils.CheckCodeWithSuite(s, err, codes.NotFound)

	s.Require().Contains(message, domain.ErrProductNotFound.Error())
}
