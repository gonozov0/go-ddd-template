package helpers

import (
	"context"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"fmt"

	pb "go-ddd-template/generated/server"
	"go-ddd-template/integration_tests/suites"
	productunithelpers "go-ddd-template/internal/application/products/shared/helpers"
	domain "go-ddd-template/internal/domain/products"
	"go-ddd-template/internal/domain/shared/valueobjects"
	"go-ddd-template/pkg/backoff"
	"go-ddd-template/pkg/grpcutils"
	loggerutils "go-ddd-template/pkg/logger/utils"
)

type HelperSuite struct {
	baseSuite *suites.BaseSuite

	grpcClient pb.ProductServiceClient
}

func NewHelperSuite(baseSuite *suites.BaseSuite) *HelperSuite {
	handlerSuite := &HelperSuite{
		baseSuite: baseSuite,

		grpcClient: pb.NewProductServiceClient(baseSuite.ServerConn),
	}

	return handlerSuite
}

func (s *HelperSuite) CreateRandomProducts(
	count int,
) valueobjects.ProductIDs {
	products := make(domain.Products, 0, count)
	for range count {
		products = append(products, *productunithelpers.GenerateProduct(s.baseSuite))
	}

	req := productunithelpers.ToCreateProductsRequest(products)

	_, err := s.grpcClient.CreateProducts(s.baseSuite.AdminCtx, req)
	grpcutils.CheckCodeWithSuite(s.baseSuite, err, codes.OK)
	s.WaitPublishProductsStatus(products.IDs())

	s.baseSuite.T().Cleanup(func() { s.DeleteProductsInCleanup(products.IDs()) })

	return products.IDs()
}

func (s *HelperSuite) CreateProducts(
	products domain.Products,
) valueobjects.ProductIDs {
	req := productunithelpers.ToCreateProductsRequest(products)

	_, err := s.grpcClient.CreateProducts(s.baseSuite.AdminCtx, req)
	grpcutils.CheckCodeWithSuite(s.baseSuite, err, codes.OK)
	s.WaitPublishProductsStatus(products.IDs())

	s.baseSuite.T().Cleanup(func() { s.DeleteProductsInCleanup(products.IDs()) })

	return products.IDs()
}

func (s *HelperSuite) DeleteProductsInCleanup(
	productIDs valueobjects.ProductIDs,
) {
	if len(productIDs) == 0 {
		return
	}

	_, err := s.grpcClient.DeleteProducts(s.baseSuite.AdminCtx, &pb.DeleteProductsRequest{
		Ids: productIDs.Strings(),
	})

	grpcStatus, ok := status.FromError(err)
	if err != nil || !ok || grpcStatus.Code() != codes.OK {
		slog.Error(
			"failed to delete products",
			loggerutils.ErrAttr(err),
			loggerutils.Attr("grpc_status", grpcStatus.Code().String()),
			loggerutils.Attr("grpc_message", grpcStatus.Message()),
			loggerutils.Attr("product_ids", productIDs.Strings()),
		)

		return
	}

	slog.Info("Deleted products", loggerutils.Attr("product_ids", productIDs.Strings()))
}

func (s *HelperSuite) WaitPublishProductsStatus(productIDs valueobjects.ProductIDs) {
	err := backoff.RunWithRetry(context.Background(), func() error {
		resp, err := s.grpcClient.GetProducts(s.baseSuite.AdminCtx, &pb.GetProductsRequest{
			Ids: productIDs.Strings(),
		})

		_, err = grpcutils.CheckCode(err, codes.OK)
		if err != nil {
			return backoff.NewAlwaysRetryableError(fmt.Errorf("failed to check code: %w", err))
		}

		if resp == nil {
			return backoff.NewAlwaysRetryableError(fmt.Errorf("reponse should not be nil"))
		}

		for _, product := range resp.Items {
			if product.Status != valueobjects.ProductStatusPublished.String() {
				return backoff.NewAlwaysRetryableError(fmt.Errorf("product %s is not published", product.Id))
			}
		}

		return nil
	}, backoff.DefaultRetryLimit)
	s.baseSuite.Require().NoError(err)
}
