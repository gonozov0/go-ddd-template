package infrastructure_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/topic/topicoptions"

	"fmt"

	"go-ddd-template/internal"
	deliveriesconsumers "go-ddd-template/internal/application/deliveries/consumers"
	productunithelpers "go-ddd-template/internal/application/products/shared/helpers"
	userunithelpers "go-ddd-template/internal/application/users/shared/helpers"
	ordersdomain "go-ddd-template/internal/domain/orders"
	productsdomain "go-ddd-template/internal/domain/products"
	"go-ddd-template/internal/domain/shared/valueobjects"
	usersdomain "go-ddd-template/internal/domain/users"
	ordersinfra "go-ddd-template/internal/infrastructure/orders"
	productsinfra "go-ddd-template/internal/infrastructure/products"
	usersinfra "go-ddd-template/internal/infrastructure/users"
	"go-ddd-template/internal/shared"
	pgmigrations "go-ddd-template/migrations/postgres"
	ydbmigrations "go-ddd-template/migrations/ydb"
	"go-ddd-template/pkg/consumerutils"
	"go-ddd-template/pkg/db/postgres"
	pgsuites "go-ddd-template/pkg/db/postgres/suites"
	"go-ddd-template/pkg/db/redis"
	redissuites "go-ddd-template/pkg/db/redis/suites"
	pkgydb "go-ddd-template/pkg/db/ydb"
	ydbsuites "go-ddd-template/pkg/db/ydb/suites"
	loggerutils "go-ddd-template/pkg/logger/utils"
	"go-ddd-template/pkg/sqs"
)

type orderCreatedEvent struct {
	OrderID uuid.UUID `json:"order_id"`
}

type DBRepoSuite struct {
	suite.Suite
	pgsuites.RunPostgresSuite
	redissuites.RunRedisSuite
	ydbsuites.RunYDBSuite

	cluster     *postgres.Cluster
	redisClient *redis.Client
	driver      *ydb.Driver
	sqsSession  sqs.Session

	orderQueueWriters   ordersinfra.QueueWriters
	orderQueueReaders   deliveriesconsumers.QueueReaders
	productQueueWriters productsinfra.QueueWriters

	ordersRepo   *ordersinfra.DBRepo
	productsRepo *productsinfra.PostgresRepo
	usersRepo    *usersinfra.DBRepo
}

func (s *DBRepoSuite) SetupSuite() {
	cfg, err := internal.LoadConfig()
	if err != nil {
		s.Fail("Failed to load config", err)
	}

	cfg.Postgres = s.RunPostgresSuite.SetupSuite(&s.Suite, s.T(), cfg.Postgres, pgmigrations.FS)
	cfg.Redis = s.RunRedisSuite.SetupSuite(&s.Suite, s.T(), cfg.Redis)
	s.RunYDBSuite.SetupSuite(&s.Suite, s.T(), cfg.YDB, ydbmigrations.FS)

	s.cluster, err = postgres.NewCluster(context.Background(), cfg.Postgres)
	if err != nil {
		s.Fail("Failed to init postgres cluster", err)
	}

	s.driver, err = pkgydb.InitNative(context.Background(), cfg.YDB)
	if err != nil {
		s.Fail("Failed to init ydb driver", err)
	}

	s.sqsSession, err = sqs.NewSession(cfg.SQS)
	if err != nil {
		s.Fail("Failed to init sqs session", err)
	}

	s.orderQueueWriters.OrderCreated, err = s.driver.Topic().StartWriter(shared.OrderProcessingTopic)
	if err != nil {
		s.Fail("Failed to start order_created topic writer", err)
	}

	s.orderQueueReaders.OrderProcessingDelivery, err = s.driver.Topic().StartReader(
		shared.OrderProcessingConsumer,
		topicoptions.ReadTopic(shared.OrderProcessingTopic),
	)
	if err != nil {
		s.Fail("Failed to start order_created topic writer", err)
	}

	s.productQueueWriters.ProductInited = sqs.NewWriter(s.sqsSession, shared.ProductInitedQueue)
	if err != nil {
		s.Fail("Failed to start product_inited topic writer", err)
	}

	s.ordersRepo = ordersinfra.NewDBRepo(s.cluster, s.orderQueueWriters)
	s.productsRepo = productsinfra.NewDBRepo(s.cluster, s.productQueueWriters)

	s.redisClient, err = redis.NewClient(context.Background(), cfg.Redis)
	if err != nil {
		s.Fail("Failed to init redis client", err)
	}

	s.usersRepo = usersinfra.NewDBRepo(s.cluster, s.redisClient)
}

func (s *DBRepoSuite) TearDownSuite() {
	if s.redisClient != nil {
		if err := s.redisClient.Close(); err != nil {
			slog.Error("Failed to close redis client", loggerutils.ErrAttr(err))
		}
	}

	if err := s.cluster.Close(); err != nil {
		s.Fail("Failed to close postgres cluster", err)
	}

	if err := s.orderQueueWriters.OrderCreated.Close(context.Background()); err != nil {
		s.Fail("Failed to close order_created topic writer", err)
	}

	if err := s.orderQueueReaders.OrderProcessingDelivery.Close(context.Background()); err != nil {
		s.Fail("Failed to close order_created topic reader", err)
	}

	if err := s.driver.Close(context.Background()); err != nil {
		s.Fail("Failed to close postgres driver", err)
	}
}

func (s *DBRepoSuite) TestOrderCreate() {
	userID, productIDs := s.prepare()

	var createdOrder *ordersdomain.Order

	s.Run("Creating order should be successful", func() {
		var err error

		createdOrder, err = s.ordersRepo.CreateOrder(
			context.Background(),
			productIDs,
			func(products []ordersdomain.Product) (*ordersdomain.Order, error) {
				for i, product := range products {
					if err := product.Reserve(); err != nil {
						return nil, fmt.Errorf("failed to reserve product: %w", err)
					}

					products[i] = product
				}

				return ordersdomain.CreateOrder(userID, products), nil
			},
		)
		s.Require().NoError(err)

		order, err := s.ordersRepo.GetOrder(context.Background(), createdOrder.GetID())
		s.Require().NoError(err)
		s.Require().Equal(createdOrder, order)
	})

	s.T().Cleanup(func() {
		if err := s.ordersRepo.DeleteOrder(context.Background(), createdOrder.GetID()); err != nil {
			slog.Error("failed to delete order", loggerutils.ErrAttr(err))
		}
	})

	s.Run("Creation order with same products should return error", func() {
		_, err := s.ordersRepo.CreateOrder(
			context.Background(),
			productIDs,
			func(products []ordersdomain.Product) (*ordersdomain.Order, error) {
				for i, product := range products {
					if err := product.Reserve(); err != nil {
						return nil, fmt.Errorf("failed to reserve product: %w", err)
					}

					products[i] = product
				}

				return ordersdomain.CreateOrder(userID, products), nil
			},
		)

		s.Require().Error(err)
		s.Require().ErrorIs(err, ordersdomain.ErrProductAlreadyReserved)
	})

	s.Run("Sending for delivery created orders should be successful", func() {
		err := s.ordersRepo.ProcessOrders(
			context.Background(),
			func(orders []*ordersdomain.Order) error {
				for _, order := range orders {
					if err := order.Process(); err != nil {
						return fmt.Errorf("failed to set processing status: %w", err)
					}
				}

				return nil
			},
		)
		s.Require().NoError(err)

		s.Require().NoError(
			consumerutils.ReadAndProcessMessageFromYDB(context.Background(), s.orderQueueReaders.OrderProcessingDelivery,
				func(event orderCreatedEvent) {
					s.Require().Equal(createdOrder.GetID().String(), event.OrderID.String())
				},
			),
		)

		actualOrder, err := s.ordersRepo.GetOrder(context.Background(), createdOrder.GetID())
		s.Require().NoError(err)

		s.Require().Equal(ordersdomain.OrderStatusProcessing, actualOrder.GetStatus())
	})
}

func (s *DBRepoSuite) TestOrderNotFound() {
	order, err := s.ordersRepo.GetOrder(context.Background(), valueobjects.EmptyOrderID)
	s.Require().Nil(order)
	s.Require().ErrorIs(err, ordersdomain.ErrOrderNotFound)
}

func (s *DBRepoSuite) prepare() (valueobjects.UserID, valueobjects.ProductIDs) {
	user := userunithelpers.GenerateUser(s)
	_, err := s.usersRepo.CreateUser(context.Background(), func() (*usersdomain.User, error) { return user, nil })
	s.Require().NoError(err)

	s.T().Cleanup(func() {
		if err = s.usersRepo.DeleteUser(context.Background(), user.GetID(), func(u *usersdomain.User) error { return nil }); err != nil {
			slog.Error("failed to delete user", loggerutils.ErrAttr(err))
		}
	})

	products := productsdomain.Products{
		*productunithelpers.GenerateProduct(s,
			productunithelpers.ProductWithStatus(valueobjects.ProductStatusPublished),
		),
		*productunithelpers.GenerateProduct(s,
			productunithelpers.ProductWithStatus(valueobjects.ProductStatusPublished),
		),
	}

	s.T().Cleanup(func() {
		if err = s.productsRepo.DeleteProducts(context.Background(), products.IDs(), func(_ []productsdomain.Product) error {
			return nil
		}); err != nil {
			slog.Error("failed to delete products", loggerutils.ErrAttr(err))
		}
	})

	err = s.productsRepo.CreateProductsInTests(
		context.Background(),
		products,
	)
	s.Require().NoError(err)

	return user.GetID(), products.IDs()
}

func TestPostgresRepoSuite(t *testing.T) {
	suite.Run(t, new(DBRepoSuite))
}
