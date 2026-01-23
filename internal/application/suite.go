package application

import (
	"context"

	"github.com/stretchr/testify/suite"

	"go-ddd-template/internal/domain/shared/valueobjects"
	"go-ddd-template/pkg/auth"
	"go-ddd-template/pkg/auth/tests"

	deliveriesinfra "go-ddd-template/internal/infrastructure/deliveries"
	ordersinfra "go-ddd-template/internal/infrastructure/orders"
	productsinfra "go-ddd-template/internal/infrastructure/products"
	usersinfra "go-ddd-template/internal/infrastructure/users"
	imagestorage "go-ddd-template/pkg/image_storage"
)

type ServerSuite struct {
	suite.Suite

	UserID   valueobjects.UserID
	UserCtx  context.Context
	AdminCtx context.Context

	UsersRepo      *usersinfra.InMemoryRepo
	OrdersRepo     *ordersinfra.InMemoryRepo
	ProductsRepo   *productsinfra.InMemoryRepo
	DeliveriesRepo *deliveriesinfra.InMemoryRepo
	ImageStorage   *imagestorage.InMemoryClient
}

func (s *ServerSuite) SetupTest() {
	s.UserID = valueobjects.UserID(tests.UserID)
	s.UserCtx = auth.WithUserInfo(context.Background(), auth.NewUserInfo(tests.UserID))
	s.AdminCtx = auth.WithUserInfo(context.Background(), auth.NewUserInfo(tests.AdminUserID))

	s.UsersRepo = usersinfra.NewInMemoryRepo()
	s.OrdersRepo = ordersinfra.NewInMemoryRepo()
	s.ProductsRepo = productsinfra.NewInMemoryRepo()
	s.DeliveriesRepo = deliveriesinfra.NewInMemoryRepo()
	s.ImageStorage = imagestorage.NewInMemoryClient(context.Background())
}
