package products_test

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"

	pb "go-ddd-template/generated/server"
	productunithelpers "go-ddd-template/internal/application/products/shared/helpers"
	domain "go-ddd-template/internal/domain/products"
	"go-ddd-template/pkg/grpcutils"
	imagestorage "go-ddd-template/pkg/image_storage"
)

func (s *ProductsSuite) TestConfirmProductImageUpload() {
	product := productunithelpers.CreateRandomProducts(s, s.ProductsRepo, 1)[0]

	_, err := s.GRPCHandlers.GetProductImageUploadURL(
		context.Background(),
		&pb.GetProductImageUploadURLRequest{
			Id:       product.GetID().String(),
			Filename: "test.jpg",
		},
	)
	s.Require().NoError(err)

	req := &pb.ConfirmProductImageUploadRequest{
		Id: product.GetID().String(),
	}
	resp, err := s.GRPCHandlers.ConfirmProductImageUpload(context.Background(), req)
	grpcutils.CheckCodeWithSuite(s, err, codes.OK)

	s.Require().NotEmpty(resp.GetImageUrl())
}

func (s *ProductsSuite) TestConfirmProductImageUploadWithInvalidProductID() {
	req := &pb.ConfirmProductImageUploadRequest{
		Id: "invalid-id",
	}
	_, err := s.GRPCHandlers.ConfirmProductImageUpload(context.Background(), req)
	message := grpcutils.CheckCodeWithSuite(s, err, codes.InvalidArgument)

	s.Require().Contains(message, "invalid product id")
}

func (s *ProductsSuite) TestConfirmProductImageUploadWithNonExistentProductID() {
	req := &pb.ConfirmProductImageUploadRequest{
		Id: uuid.New().String(),
	}
	_, err := s.GRPCHandlers.ConfirmProductImageUpload(context.Background(), req)
	message := grpcutils.CheckCodeWithSuite(s, err, codes.NotFound)

	s.Require().Contains(message, domain.ErrProductNotFound.Error())
}

func (s *ProductsSuite) TestConfirmProductImageUploadWithImageNotFound() {
	product := productunithelpers.CreateRandomProducts(s, s.ProductsRepo, 1)[0]

	req := &pb.ConfirmProductImageUploadRequest{
		Id: product.GetID().String(),
	}
	_, err := s.GRPCHandlers.ConfirmProductImageUpload(context.Background(), req)
	message := grpcutils.CheckCodeWithSuite(s, err, codes.FailedPrecondition)

	s.Require().Contains(message, imagestorage.ErrImageNotFound.Error())
}
