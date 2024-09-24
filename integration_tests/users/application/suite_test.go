package application_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"go-echo-template/generated/openapi"
	"go-echo-template/generated/protobuf"
	"go-echo-template/integration_tests/suites"
	"go-echo-template/internal"
	infra "go-echo-template/internal/infrastructure/users"
	"go-echo-template/pkg/postgres"
)

type UsersSuite struct {
	suite.Suite
	suites.RunServerSuite

	grpcClient protobuf.UserServiceClient
	httpClient *openapi.Client
	repo       *infra.PostgresRepo
}

func (s *UsersSuite) SetupSuite() {
	s.RunServerSuite.SetupSuite("8082")
	s.grpcClient = protobuf.NewUserServiceClient(s.Conn)
	var err error
	s.httpClient, err = openapi.NewClient(s.HTTPServerURL)
	if err != nil {
		s.Fail("Failed to create HTTP client", err)
	}
	s.repo, err = initPostgresRepo()
	if err != nil {
		s.Fail("Failed to init postgres repo", err)
	}
}

func (s *UsersSuite) TearDownSuite() {
	s.RunServerSuite.TearDownSuite()
}

func TestUsersSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(UsersSuite))
}

func initPostgresRepo() (*infra.PostgresRepo, error) {
	cfg, err := internal.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
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
		return nil, fmt.Errorf("failed to init postgres connection data: %w", err)
	}
	cluster, err := postgres.InitCluster(context.Background(), connData)
	if err != nil {
		return nil, fmt.Errorf("failed to init postgres cluster: %w", err)
	}
	return infra.NewPostgresRepo(cluster), nil
}
