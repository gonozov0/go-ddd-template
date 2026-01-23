package products_test

import (
	"context"

	"google.golang.org/grpc/codes"

	productunithelpers "go-ddd-template/internal/application/products/shared/helpers"
	domain "go-ddd-template/internal/domain/products"
	"go-ddd-template/pkg/grpcutils"
)

func (s *ProductsSuite) TestCreateProducts() {
	products := domain.Products{
		*productunithelpers.GenerateProduct(s),
		*productunithelpers.GenerateProduct(s),
	}

	resp, err := s.GRPCHandlers.CreateProducts(
		context.Background(),
		productunithelpers.ToCreateProductsRequest(products),
	)
	grpcutils.CheckCodeWithSuite(s, err, codes.OK)

	s.Require().Equal(2, len(resp.GetIds()))
}

func (s *ProductsSuite) TestCreateProductsWithErrInvalidProductName() {
	products := domain.Products{
		*productunithelpers.GenerateProduct(s),
		*productunithelpers.GenerateProduct(s,
			productunithelpers.ProductWithName(""),
		),
	}
	_, err := s.GRPCHandlers.CreateProducts(context.Background(), productunithelpers.ToCreateProductsRequest(products))
	message := grpcutils.CheckCodeWithSuite(s, err, codes.InvalidArgument)

	s.Require().Contains(message, domain.ErrInvalidProduct.Error())
	s.Require().Contains(message, "name must not be empty")
}

func (s *ProductsSuite) TestCreateProductsWithErrInvalidProductPrice() {
	products := domain.Products{
		*productunithelpers.GenerateProduct(s),
		*productunithelpers.GenerateProduct(s,
			productunithelpers.ProductWithPrice(0),
		),
	}
	_, err := s.GRPCHandlers.CreateProducts(context.Background(), productunithelpers.ToCreateProductsRequest(products))
	message := grpcutils.CheckCodeWithSuite(s, err, codes.InvalidArgument)

	s.Require().Contains(message, domain.ErrInvalidProduct.Error())
	s.Require().Contains(message, "price must be greater than 0")
}
