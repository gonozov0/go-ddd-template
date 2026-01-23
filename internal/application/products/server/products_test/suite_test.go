package products_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"go-ddd-template/internal/application"
	products "go-ddd-template/internal/application/products/server"
)

type ProductsSuite struct {
	suite.Suite
	application.ServerSuite
	GRPCHandlers products.ProductHandlers
}

func (s *ProductsSuite) SetupTest() {
	s.ServerSuite.SetupTest()
	s.GRPCHandlers = products.SetupHandlers(s.ProductsRepo, s.ImageStorage)
}

func TestProductsSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ProductsSuite))
}
