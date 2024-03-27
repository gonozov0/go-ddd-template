package users_test

import (
	"testing"

	"go-echo-template/internal/infrastructure/users"

	"github.com/stretchr/testify/suite"
)

type PostgresRepoSuite struct {
	suite.Suite
	repo *users.PostgresRepo
}

func (suite *PostgresRepoSuite) SetupSuite() {
	suite.repo = users.NewPostgresRepo()
}

func (suite *PostgresRepoSuite) TestSaveUser() {
	suite.T().Skip("not implemented")
}

func TestPostgresRepoSuite(t *testing.T) {
	suite.Run(t, new(PostgresRepoSuite))
}
