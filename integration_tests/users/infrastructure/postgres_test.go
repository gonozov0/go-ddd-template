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
	email := "test@test.com"
	created, err := suite.repo.CreateUser(context.Background(), email, func() (*domain.User, error) {
		return domain.CreateUser("test", email)
	})
	suite.Require().NoError(err)

	gotten, err := suite.repo.GetUser(context.Background(), created.ID())
	suite.Require().NoError(err)
	suite.Require().Equal(created, gotten)

	updated, err := suite.repo.UpdateUser(context.Background(), created.ID(), func(u *domain.User) (bool, error) {
		err := u.ChangeEmail("test@test2.com")
		return true, err
	})
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
