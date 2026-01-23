package infrastructure_test

import (
	"context"
	"testing"

	"go-ddd-template/pkg/slices"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/suite"

	"fmt"

	"go-ddd-template/internal"
	productconsumers "go-ddd-template/internal/application/products/consumers"
	productunithelpers "go-ddd-template/internal/application/products/shared/helpers"
	domain "go-ddd-template/internal/domain/products"
	"go-ddd-template/internal/domain/shared/valueobjects"
	infra "go-ddd-template/internal/infrastructure/products"
	"go-ddd-template/internal/shared"
	pgmigrations "go-ddd-template/migrations/postgres"
	"go-ddd-template/pkg/consumerutils"
	"go-ddd-template/pkg/db/postgres"
	pgsuites "go-ddd-template/pkg/db/postgres/suites"
	"go-ddd-template/pkg/sqs"
)

type PostgresRepoSuite struct {
	suite.Suite

	pgsuites.RunPostgresSuite

	cluster    *postgres.Cluster
	sqsSession sqs.Session

	productQueueWriters infra.QueueWriters
	productQueueReaders productconsumers.QueueReaders

	repo *infra.PostgresRepo
}

func (s *PostgresRepoSuite) SetupSuite() {
	s.RunPostgresSuite.SetS(&s.Suite)
	s.RunPostgresSuite.SetT(s.T())

	cfg, err := internal.LoadConfig()
	if err != nil {
		s.Fail("Failed to load config", err)
	}

	cfg.Postgres = s.RunPostgresSuite.SetupSuite(&s.Suite, s.T(), cfg.Postgres, pgmigrations.FS)

	s.cluster, err = postgres.NewCluster(context.Background(), cfg.Postgres)
	if err != nil {
		s.Fail("Failed to init postgres cluster", err)
	}

	s.sqsSession, err = sqs.NewSession(cfg.SQS)
	if err != nil {
		s.Fail("Failed to init sqs session", err)
	}

	s.productQueueWriters.ProductInited = sqs.NewWriter(s.sqsSession, shared.ProductInitedQueue)
	if err != nil {
		s.Fail("Failed to start product_inited topic writer", err)
	}

	s.productQueueReaders.ProductInited = sqs.NewReader(
		s.sqsSession,
		shared.ProductInitedQueue,
		sqs.WithMaxNumberOfMessages(1),
	)

	s.repo = infra.NewDBRepo(s.cluster, s.productQueueWriters)
}

func (s *PostgresRepoSuite) TearDownSuite() {
	if err := s.cluster.Close(); err != nil {
		s.Fail("Failed to close postgres cluster", err)
	}
}

func (s *PostgresRepoSuite) TestProductCRUD() {
	created := domain.Products{
		*productunithelpers.GenerateProduct(s),
		*productunithelpers.GenerateProduct(s),
		*productunithelpers.GenerateProduct(s),
	}

	createdIDs := slices.Map(created, func(p domain.Product) valueobjects.ProductID {
		return p.GetID()
	})

	err := s.repo.CreateProducts(context.Background(), func() (domain.Products, error) {
		return created, nil
	})
	s.Require().NoError(err)

	for range len(created) {
		s.Require().
			NoError(consumerutils.ReadAndProcessMessageFromSQS(context.Background(), s.productQueueReaders.ProductInited, func(event ProductInitedEvent) error {
				if !slices.Contains(created.IDs().Strings(), event.ID) {
					return fmt.Errorf("product with id %s not found", event.ID)
				}

				return nil
			}))
	}

	s.T().Cleanup(func() {
		if err = s.repo.DeleteProducts(context.Background(), createdIDs, func(_ []domain.Product) error {
			return nil
		}); err != nil {
			s.T().Error(err)
		}
	})

	gotten, err := s.repo.GetProducts(context.Background(), createdIDs)
	s.Require().NoError(err)

	s.Require().ElementsMatch(created, gotten)

	updatedImageURL := valueobjects.NewImageURL(gofakeit.URL())

	err = s.repo.UpdateProduct(context.Background(), createdIDs[0], func(product *domain.Product) error {
		product.SetImageURL(updatedImageURL)
		return nil
	})
	s.Require().NoError(err)

	updated, err := s.repo.GetProducts(context.Background(), []valueobjects.ProductID{createdIDs[0]})
	s.Require().NoError(err)
	s.Require().Len(updated, 1)
	s.Require().Equal(updatedImageURL, updated[0].GetImageURL())
}

func TestPostgresRepoSuite(t *testing.T) {
	suite.Run(t, new(PostgresRepoSuite))
}

type ProductInitedEvent struct {
	ID string `json:"id"`
}
