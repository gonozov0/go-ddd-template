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

func (s *ProductsSuite) TestGetProductImageUploadURL() {
	product := productunithelpers.CreateRandomProducts(s, s.ProductsRepo, 1)[0]

	req := &pb.GetProductImageUploadURLRequest{
		Id:       product.GetID().String(),
		Filename: "test.jpg",
	}
	resp, err := s.GRPCHandlers.GetProductImageUploadURL(context.Background(), req)
	grpcutils.CheckCodeWithSuite(s, err, codes.OK)

	s.Require().NotEmpty(resp.GetUploadUrl())
}

func (s *ProductsSuite) TestGetProductImageUploadURLWithInvalidProductID() {
	req := &pb.GetProductImageUploadURLRequest{
		Id:       "invalid-id",
		Filename: "test.jpg",
	}
	_, err := s.GRPCHandlers.GetProductImageUploadURL(context.Background(), req)
	message := grpcutils.CheckCodeWithSuite(s, err, codes.InvalidArgument)

	s.Require().Contains(message, "invalid product id")
}

func (s *ProductsSuite) TestGetProductImageUploadURLWithInvalidFilename() {
	product := productunithelpers.CreateRandomProducts(s, s.ProductsRepo, 1)[0]

	req := &pb.GetProductImageUploadURLRequest{
		Id:       product.GetID().String(),
		Filename: "",
	}
	_, err := s.GRPCHandlers.GetProductImageUploadURL(context.Background(), req)
	message := grpcutils.CheckCodeWithSuite(s, err, codes.InvalidArgument)

	s.Require().Contains(message, "invalid image filename")
	s.Require().Contains(message, "filename must not be empty")
}

func (s *ProductsSuite) TestGetProductImageUploadURLWithNonExistentProductID() {
	req := &pb.GetProductImageUploadURLRequest{
		Id:       uuid.New().String(),
		Filename: "test.jpg",
	}
	_, err := s.GRPCHandlers.GetProductImageUploadURL(context.Background(), req)
	message := grpcutils.CheckCodeWithSuite(s, err, codes.NotFound)

	s.Require().Contains(message, domain.ErrProductNotFound.Error())
}
