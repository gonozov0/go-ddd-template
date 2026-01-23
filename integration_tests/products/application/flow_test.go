package application_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"google.golang.org/grpc/codes"

	pb "go-ddd-template/generated/server"
	productunithelpers "go-ddd-template/internal/application/products/shared/helpers"
	domain "go-ddd-template/internal/domain/products"
	"go-ddd-template/internal/domain/shared/valueobjects"
	"go-ddd-template/pkg/grpcutils"
	imagestorage "go-ddd-template/pkg/image_storage"
	"go-ddd-template/pkg/slices"
)

type CompareItem struct {
	Name  string
	Price float64
	ID    string
}

func (s *ProductsSuite) TestGRPC() {
	var (
		products = domain.Products{
			*productunithelpers.GenerateProduct(&s.BaseSuite,
				productunithelpers.ProductWithStatus(valueobjects.ProductStatusInit),
			),
			*productunithelpers.GenerateProduct(&s.BaseSuite,
				productunithelpers.ProductWithStatus(valueobjects.ProductStatusInit),
			),
			*productunithelpers.GenerateProduct(&s.BaseSuite,
				productunithelpers.ProductWithStatus(valueobjects.ProductStatusInit),
			),
		}
	)

	s.Run("Create Products", func() {
		resp, err := s.grpcClient.CreateProducts(s.AdminCtx, productunithelpers.ToCreateProductsRequest(products))
		grpcutils.CheckCodeWithSuite(s, err, codes.OK)
		s.Require().NotNil(resp)

		s.Require().NotEmpty(resp.GetIds())

		productIDs, err := valueobjects.NewProductIDsFromStrings(resp.GetIds())
		s.Require().NoError(err)

		products, err = productunithelpers.UpdateProductsWithIDs(productIDs, products)
		s.Require().NoError(err)
	})

	s.T().Cleanup(func() {
		s.productHelper.DeleteProductsInCleanup(products.IDs())
	})

	products = slices.Map(products, func(p domain.Product) domain.Product {
		s.Require().NoError(p.Publish())
		return p
	})

	s.productHelper.WaitPublishProductsStatus(products.IDs())

	s.Run("Get Products", func() {
		s.Require().NotEmpty(products.IDs())

		getReq := &pb.GetProductsRequest{
			Ids: products.IDs().Strings(),
		}

		resp, err := s.grpcClient.GetProducts(s.AdminCtx, getReq)
		grpcutils.CheckCodeWithSuite(s, err, codes.OK)
		s.Require().NotNil(resp)

		expectedProducts := make([]*pb.ProductResponseItem, 0, len(products))
		for _, product := range products {
			expectedProducts = append(expectedProducts, &pb.ProductResponseItem{
				Id:     product.GetID().String(),
				Name:   product.GetName().String(),
				Price:  product.GetPrice().Float64(),
				Status: product.GetStatus().String(),
			})
		}

		s.Require().ElementsMatch(expectedProducts, resp.Items)
	})

	s.Run("Product Image Upload Flow", func() {
		s.Require().NotEmpty(products)

		productID := products[0].GetID().String()
		filename := "test-image.jpg"
		imageContent := []byte("integration test image content")

		uploadURLResp, err := s.grpcClient.GetProductImageUploadURL(
			s.AdminCtx,
			&pb.GetProductImageUploadURLRequest{
				Id:       productID,
				Filename: filename,
			},
		)
		grpcutils.CheckCodeWithSuite(s, err, codes.OK)
		s.Require().NotNil(uploadURLResp)
		s.Require().NotEmpty(uploadURLResp.GetUploadUrl())

		uploadReq, err := http.NewRequest(http.MethodPut, uploadURLResp.GetUploadUrl(), bytes.NewReader(imageContent))
		s.Require().NoError(err)

		uploadReq.Header.Set("x-amz-acl", imagestorage.ImageBucketACL)

		uploadResp, err := http.DefaultClient.Do(uploadReq)
		s.Require().NoError(err)

		_, err = io.Copy(io.Discard, uploadResp.Body)
		s.Require().NoError(err)
		s.Require().NoError(uploadResp.Body.Close())
		s.Require().True(
			uploadResp.StatusCode == http.StatusOK || uploadResp.StatusCode == http.StatusNoContent,
			"unexpected status code from upload url: %d",
			uploadResp.StatusCode,
		)

		confirmResp, err := s.grpcClient.ConfirmProductImageUpload(
			s.AdminCtx,
			&pb.ConfirmProductImageUploadRequest{
				Id: productID,
			},
		)
		grpcutils.CheckCodeWithSuite(s, err, codes.OK)
		s.Require().NotNil(confirmResp)
		s.Require().NotEmpty(confirmResp.GetImageUrl())

		publicResp, err := http.Get(confirmResp.GetImageUrl())
		s.Require().NoError(err)

		downloadedContent, err := io.ReadAll(publicResp.Body)
		s.Require().NoError(err)
		s.Require().NoError(publicResp.Body.Close())
		s.Require().Equal(http.StatusOK, publicResp.StatusCode, string(downloadedContent))
		s.Require().Equal(imageContent, downloadedContent)
	})
}

func (s *ProductsSuite) TestHTTP() {
	var (
		products = domain.Products{
			*productunithelpers.GenerateProduct(&s.BaseSuite,
				productunithelpers.ProductWithStatus(valueobjects.ProductStatusInit),
			),
			*productunithelpers.GenerateProduct(&s.BaseSuite,
				productunithelpers.ProductWithStatus(valueobjects.ProductStatusInit),
			),
			*productunithelpers.GenerateProduct(&s.BaseSuite,
				productunithelpers.ProductWithStatus(valueobjects.ProductStatusInit),
			),
		}
	)

	s.Run("Create Products", func() {
		body, err := grpcutils.MarshalJSON(productunithelpers.ToCreateProductsRequest(products))
		s.Require().NoError(err)

		httpReq, err := http.NewRequest(
			http.MethodPost,
			fmt.Sprintf("%s/products", s.ServerURL),
			bytes.NewBuffer(body),
		)
		s.Require().NoError(err)

		httpReq.Header = s.AdminHeaders

		httpResp, err := http.DefaultClient.Do(httpReq)
		s.Require().NoError(err)

		respBytes, err := io.ReadAll(httpResp.Body)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, httpResp.StatusCode, string(respBytes))

		var resp pb.CreateProductsResponse

		err = grpcutils.UnmarshalJSON(respBytes, &resp)
		s.Require().NoError(err)
		s.Require().NoError(httpResp.Body.Close())

		productIDs, err := valueobjects.NewProductIDsFromStrings(resp.GetIds())
		s.Require().NoError(err)

		products, err = productunithelpers.UpdateProductsWithIDs(productIDs, products)
		s.Require().NoError(err)
	})

	s.T().Cleanup(func() {
		s.productHelper.DeleteProductsInCleanup(products.IDs())
	})

	products = slices.Map(products, func(p domain.Product) domain.Product {
		s.Require().NoError(p.Publish())
		return p
	})

	s.productHelper.WaitPublishProductsStatus(products.IDs())

	s.Run("Get Products", func() {
		s.Require().NotEmpty(products.IDs())

		body, err := grpcutils.MarshalJSON(products.IDs().Strings())
		s.Require().NoError(err)

		httpReq, err := http.NewRequest(
			http.MethodPost,
			fmt.Sprintf("%s/products/get_list", s.ServerURL),
			bytes.NewBuffer(body),
		)
		s.Require().NoError(err)

		httpReq.Header = s.AdminHeaders

		httpResp, err := http.DefaultClient.Do(httpReq)
		s.Require().NoError(err)

		respBytes, err := io.ReadAll(httpResp.Body)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, httpResp.StatusCode, string(respBytes))

		var resp *pb.GetProductsResponse

		err = grpcutils.UnmarshalJSON(respBytes, &resp)
		s.Require().NoError(err)
		s.Require().NoError(httpResp.Body.Close())

		expectedProducts := make([]*pb.ProductResponseItem, 0, len(products))
		for _, product := range products {
			expectedProducts = append(expectedProducts, &pb.ProductResponseItem{
				Id:     product.GetID().String(),
				Name:   product.GetName().String(),
				Price:  product.GetPrice().Float64(),
				Status: product.GetStatus().String(),
			})
		}

		s.Require().ElementsMatch(expectedProducts, resp.Items)
	})

	s.Run("Product Image Upload Flow", func() {
		s.Require().NotEmpty(products)

		productID := products[0].GetID().String()
		filename := "test-image.jpg"
		imageContent := []byte("integration test image content")

		body, err := grpcutils.MarshalJSON(filename)
		s.Require().NoError(err)

		getUploadURLReq, err := http.NewRequest(
			http.MethodPost,
			fmt.Sprintf("%s/products/%s/image/upload", s.ServerURL, productID),
			bytes.NewReader(body),
		)
		s.Require().NoError(err)

		getUploadURLReq.Header = s.AdminHeaders.Clone()
		getUploadURLReq.Header.Set("Content-Type", "application/json")

		getUploadURLResp, err := http.DefaultClient.Do(getUploadURLReq)
		s.Require().NoError(err)

		respBytes, err := io.ReadAll(getUploadURLResp.Body)
		s.Require().NoError(err)
		s.Require().NoError(getUploadURLResp.Body.Close())
		s.Require().Equal(http.StatusOK, getUploadURLResp.StatusCode, string(respBytes))

		var uploadURLResp pb.GetProductImageUploadURLResponse

		err = grpcutils.UnmarshalJSON(respBytes, &uploadURLResp)
		s.Require().NoError(err)
		s.Require().NotEmpty(uploadURLResp.GetUploadUrl())

		uploadReq, err := http.NewRequest(http.MethodPut, uploadURLResp.GetUploadUrl(), bytes.NewReader(imageContent))
		s.Require().NoError(err)

		uploadReq.Header.Set("x-amz-acl", imagestorage.ImageBucketACL)

		uploadResp, err := http.DefaultClient.Do(uploadReq)
		s.Require().NoError(err)

		_, err = io.Copy(io.Discard, uploadResp.Body)
		s.Require().NoError(err)
		s.Require().NoError(uploadResp.Body.Close())
		s.Require().True(
			uploadResp.StatusCode == http.StatusOK || uploadResp.StatusCode == http.StatusNoContent,
			"unexpected status code from upload url: %d",
			uploadResp.StatusCode,
		)

		confirmReq, err := http.NewRequest(
			http.MethodPost,
			fmt.Sprintf("%s/products/%s/image/confirm", s.ServerURL, productID),
			nil,
		)
		s.Require().NoError(err)

		confirmReq.Header = s.AdminHeaders

		confirmResp, err := http.DefaultClient.Do(confirmReq)
		s.Require().NoError(err)

		respBytes, err = io.ReadAll(confirmResp.Body)
		s.Require().NoError(err)
		s.Require().NoError(confirmResp.Body.Close())
		s.Require().Equal(http.StatusOK, confirmResp.StatusCode, string(respBytes))

		var confirmURLResp pb.ConfirmProductImageUploadResponse

		err = grpcutils.UnmarshalJSON(respBytes, &confirmURLResp)
		s.Require().NoError(err)
		s.Require().NotEmpty(confirmURLResp.GetImageUrl())

		publicResp, err := http.Get(confirmURLResp.GetImageUrl())
		s.Require().NoError(err)

		downloadedContent, err := io.ReadAll(publicResp.Body)
		s.Require().NoError(err)
		s.Require().NoError(publicResp.Body.Close())
		s.Require().Equal(http.StatusOK, publicResp.StatusCode, string(downloadedContent))
		s.Require().Equal(imageContent, downloadedContent)
	})
}
