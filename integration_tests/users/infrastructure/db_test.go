package infrastructure_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/suite"

	"go-ddd-template/internal"
	userunithelpers "go-ddd-template/internal/application/users/shared/helpers"
	"go-ddd-template/internal/domain/shared/valueobjects"
	domain "go-ddd-template/internal/domain/users"
	usersinfra "go-ddd-template/internal/infrastructure/users"
	pgmigrations "go-ddd-template/migrations/postgres"
	"go-ddd-template/pkg/db/postgres"
	pgsuites "go-ddd-template/pkg/db/postgres/suites"
	"go-ddd-template/pkg/db/redis"
	redissuites "go-ddd-template/pkg/db/redis/suites"
	loggerutils "go-ddd-template/pkg/logger/utils"
)

type DBRepoSuite struct {
	suite.Suite

	pgsuites.RunPostgresSuite
	redissuites.RunRedisSuite

	postgresCluster *postgres.Cluster
	redisClient     *redis.Client
	repo            *usersinfra.DBRepo
}

func (s *DBRepoSuite) SetupSuite() {
	cfg, err := internal.LoadConfig()
	if err != nil {
		s.Fail("Failed to load config", err)
	}

	cfg.Postgres = s.RunPostgresSuite.SetupSuite(&s.Suite, s.T(), cfg.Postgres, pgmigrations.FS)

	s.postgresCluster, err = postgres.NewCluster(context.Background(), cfg.Postgres)
	if err != nil {
		s.Fail("Failed to init postgres cluster", err)
	}

	cfgRedis := s.RunRedisSuite.SetupSuite(&s.Suite, s.T(), cfg.Redis)

	s.redisClient, err = redis.NewClient(context.Background(), cfgRedis)
	if err != nil {
		s.Fail("Failed to init redis client", err)
	}

	s.repo = usersinfra.NewDBRepo(s.postgresCluster, s.redisClient)
}

func (s *DBRepoSuite) TearDownSuite() {
	if s.postgresCluster != nil {
		if err := s.postgresCluster.Close(); err != nil {
			slog.Error("failed to close postgres cluster", loggerutils.ErrAttr(err))
		}
	}

	if s.redisClient != nil {
		if err := s.redisClient.Close(); err != nil {
			slog.Error("failed to close redis client", loggerutils.ErrAttr(err))
		}
	}
}

func (s *DBRepoSuite) checkUserInBothRepos(ctx context.Context, userID valueobjects.UserID) (bool, bool) {
	var inPostgres, inRedis bool

	db, err := s.postgresCluster.StandbyPreferredDBx(ctx)
	if err == nil {
		var exists bool

		query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`

		err = db.GetContext(ctx, &exists, query, userID.String())
		if err == nil {
			inPostgres = exists
		}
	}

	_, err = s.redisClient.Get(ctx, userID.String()).Result()
	inRedis = err == nil

	return inPostgres, inRedis
}

func (s *DBRepoSuite) TestCacheHit() {
	ctx := context.Background()
	user := userunithelpers.CreateUser(&s.Suite, s.repo)

	userInPostgres, userInRedis := s.checkUserInBothRepos(ctx, user.GetID())
	s.True(userInPostgres, "User should be in postgres")
	s.True(userInRedis, "User should be in redis")

	responseUser, err := s.repo.GetUser(ctx, user.GetID())
	s.Require().NoError(err)
	s.Equal(user.GetID(), responseUser.GetID())

	responseCachedUser, err := s.repo.GetUser(ctx, user.GetID())
	s.Require().NoError(err)
	s.Equal(responseUser.GetID(), responseCachedUser.GetID())
	s.Equal(responseUser.GetName(), responseCachedUser.GetName())
	s.Equal(responseUser.GetEmail(), responseCachedUser.GetEmail())
}

func (s *DBRepoSuite) TestCacheInvalidationOnUpdate() {
	ctx := context.Background()
	user := userunithelpers.CreateUser(&s.Suite, s.repo)

	userInPostgres, userInRedis := s.checkUserInBothRepos(ctx, user.GetID())
	s.True(userInPostgres, "User should be in postgres")
	s.True(userInRedis, "User should be in redis")

	newEmail := valueobjects.Email(gofakeit.Email())
	_, err := s.repo.UpdateUser(ctx, user.GetID(), func(u *domain.User) error {
		u.Update(u.GetName(), newEmail)
		return nil
	})
	s.Require().NoError(err)

	userInPostgres, userInRedis = s.checkUserInBothRepos(ctx, user.GetID())
	s.True(userInPostgres, "User should be in postgres")
	s.True(userInRedis, "User should be in redis")

	updatedUser, err := s.repo.GetUser(ctx, user.GetID())
	s.Require().NoError(err)
	s.Equal(newEmail, updatedUser.GetEmail(), "Expecting new email from database")
}

func (s *DBRepoSuite) TestCacheInvalidationOnDelete() {
	ctx := context.Background()
	user := userunithelpers.CreateUser(&s.Suite, s.repo)

	userInPostgres, userInRedis := s.checkUserInBothRepos(ctx, user.GetID())
	s.True(userInPostgres, "User should be in postgres")
	s.True(userInRedis, "User should be in redis")

	err := s.repo.DeleteUser(ctx, user.GetID(), func(u *domain.User) error { return nil })
	s.Require().NoError(err)

	_, err = s.repo.GetUser(ctx, user.GetID())
	s.Require().Error(err)
	s.Require().ErrorIs(err, domain.ErrUserNotFound, "User should be deleted both from postgres and cache")
}

func (s *DBRepoSuite) TestUserNotFound() {
	ctx := context.Background()

	unexpectedID := valueobjects.EmptyUserID

	_, err := s.repo.GetUser(ctx, unexpectedID)
	s.Require().Error(err)
	s.Require().ErrorIs(err, domain.ErrUserNotFound)
}

func TestDBRepoSuite(t *testing.T) {
	suite.Run(t, new(DBRepoSuite))
}
