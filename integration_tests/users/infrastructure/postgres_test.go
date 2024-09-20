package infrastructure_test

import (
	"context"
	"testing"

	"go-echo-template/internal"
	domain "go-echo-template/internal/domain/users"
	"go-echo-template/pkg/postgres"

	infra "go-echo-template/internal/infrastructure/users"

	"github.com/stretchr/testify/suite"
)

type PostgresRepoSuite struct {
	suite.Suite
	repo *infra.PostgresRepo
}

func (suite *PostgresRepoSuite) SetupSuite() {
	cfg, err := internal.LoadConfig()
	if err != nil {
		suite.Fail("Failed to load config", err)
	}
	connData, err := postgres.NewConnectionData(
		cfg.Postgres.Hosts,
		cfg.Postgres.Database,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Port,
		cfg.Postgres.SSL,
	)
	if err != nil {
		suite.Fail("Failed to init postgres connection data", err)
	}
	cluster, err := postgres.InitCluster(context.Background(), connData)
	if err != nil {
		suite.Fail("Failed to init postgres cluster", err)
	}
	suite.repo = infra.NewPostgresRepo(cluster)
}

func (suite *PostgresRepoSuite) TestUserCRUD() {
	created, err := domain.CreateUser("test", "test@test.com")
	suite.Require().NoError(err)

	err = suite.repo.SaveUser(context.Background(), *created)
	suite.Require().NoError(err)

	gotten, err := suite.repo.GetUser(context.Background(), created.ID())
	suite.Require().NoError(err)
	suite.Require().Equal(created, gotten)

	updated, err := domain.NewUser(created.ID(), "test2", "test@test2.com")
	suite.Require().NoError(err)

	err = suite.repo.SaveUser(context.Background(), *updated)
	suite.Require().NoError(err)

	gotten, err = suite.repo.GetUser(context.Background(), created.ID())
	suite.Require().NoError(err)
	suite.Require().Equal(updated, gotten)

	err = suite.repo.DeleteUser(context.Background(), created.ID())
	suite.Require().NoError(err)
}

func TestPostgresRepoSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(PostgresRepoSuite))
}
