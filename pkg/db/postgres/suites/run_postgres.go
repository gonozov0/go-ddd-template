package suites

import (
	"context"
	"embed"
	"testing"

	"github.com/stretchr/testify/suite"

	"go-ddd-template/pkg/db/postgres"
)

type RunPostgresSuite struct {
	suite.Suite
}

func (s *RunPostgresSuite) SetupSuite(
	outS *suite.Suite,
	t *testing.T,
	cfg postgres.Config,
	fs embed.FS,
) postgres.Config {
	s.SetS(outS)
	s.SetT(t)

	if err := postgres.Migrate(context.Background(), cfg, fs, nil); err != nil {
		s.FailNow("Failed to migrate Postgres", err)
	}

	return cfg
}
