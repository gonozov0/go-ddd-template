package suites

import (
	"context"
	"embed"
	"testing"

	"github.com/stretchr/testify/suite"

	pkgydb "go-ddd-template/pkg/db/ydb"
)

type RunYDBSuite struct {
	suite.Suite
}

func (s *RunYDBSuite) SetupSuite(outS *suite.Suite, t *testing.T, cfg pkgydb.Config, fs embed.FS) {
	// Необходимо проинициализировать suite'ы, встроенные в каждую из структур,
	// иначе, встроенный в них suite останется с nil полями, из-за чего
	// возникает паника при вызове s.Require()  внутри suite
	s.SetS(outS)
	s.SetT(t)

	if err := pkgydb.Migrate(context.Background(), cfg, fs, nil); err != nil {
		s.FailNow("Failed to migrate YDB", err)
	}
}
