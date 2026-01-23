package application_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	pb "go-ddd-template/generated/server"
	producthelper "go-ddd-template/integration_tests/products/helpers"
	"go-ddd-template/integration_tests/suites"
)

type ProductsSuite struct {
	suites.BaseSuite

	grpcClient pb.ProductServiceClient

	productHelper *producthelper.HelperSuite
}

func (s *ProductsSuite) SetupSuite() {
	s.BaseSuite.SetupSuite()

	s.grpcClient = pb.NewProductServiceClient(s.ServerConn)

	s.productHelper = producthelper.NewHelperSuite(&s.BaseSuite)
}

func TestProductsSuite(t *testing.T) {
	suite.Run(t, new(ProductsSuite))
}
