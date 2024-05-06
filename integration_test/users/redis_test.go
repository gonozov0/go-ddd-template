package users_test

import (
	"context"
	"testing"

	"go-echo-ddd-template/internal"
	"go-echo-ddd-template/internal/domain/users"
	usersInfra "go-echo-ddd-template/internal/infrastructure/users"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type RedisRepoSuite struct {
	suite.Suite
	repo         *usersInfra.RedisRepo
	ctx          context.Context
	keysToDelete []uuid.UUID
}

func (suite *RedisRepoSuite) SetupSuite() {
	config, err := internal.LoadConfig()
	suite.Require().NoError(err)

	repo := usersInfra.NewRedisRepo(
		config.Redis.ClusterMode,
		config.Redis.TLSEnabled,
		config.Redis.Address,
		config.Redis.Username,
		config.Redis.Password,
		config.Redis.Expiration,
	)
	suite.repo = repo
	suite.ctx = context.Background()
	suite.keysToDelete = make([]uuid.UUID, 0)
}

func (suite *RedisRepoSuite) TearDownTest() {
	for _, key := range suite.keysToDelete {
		err := suite.repo.DeleteUser(suite.ctx, key)
		suite.Require().NoError(err)
	}
	suite.keysToDelete = nil
}

func (suite *RedisRepoSuite) TestSaveUser() {
	key := uuid.New()
	suite.keysToDelete = append(suite.keysToDelete, key)

	user, err := users.NewUser(key, "test", "test@test.com")
	suite.Require().NoError(err)

	err = suite.repo.SaveUser(suite.ctx, *user)
	suite.Require().NoError(err)
}

func (suite *RedisRepoSuite) TestGetUser() {
	key := uuid.New()
	suite.keysToDelete = append(suite.keysToDelete, key)

	user, err := users.NewUser(key, "test", "test@test.com")
	suite.Require().NoError(err)

	err = suite.repo.SaveUser(suite.ctx, *user)
	suite.Require().NoError(err)

	u, err := suite.repo.GetUser(suite.ctx, key)
	suite.Require().NoError(err)
	suite.Equal(user, u)
}

func (suite *RedisRepoSuite) TestGetUserNotFound() {
	_, err := suite.repo.GetUser(suite.ctx, uuid.New())
	suite.Require().ErrorIs(err, users.ErrUserNotFound)
}

func TestRedisRepoSuite(t *testing.T) {
	suite.Run(t, new(RedisRepoSuite))
}
