package suites

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	pkgredis "go-ddd-template/pkg/db/redis"
)

type RunRedisSuite struct {
	suite.Suite
}

func (s *RunRedisSuite) SetupSuite(outS *suite.Suite, t *testing.T, config pkgredis.Config) pkgredis.Config {
	// Необходимо проинициализировать suite'ы, встроенные в каждую из структур,
	// иначе, встроенный в них suite останется с nil полями, из-за чего
	// возникает паника при вызове s.Require() внутри suite
	s.SetS(outS)
	s.SetT(t)

	ports := os.Getenv("TESTSUITE_REDIS_CLUSTER_PORTS")

	if ports != "" {
		redisPorts := strings.Split(ports, ",")
		redisAddrs := make([]string, 0, len(redisPorts))

		for _, port := range redisPorts {
			port = strings.TrimSpace(port)
			if port != "" {
				redisAddrs = append(redisAddrs, fmt.Sprintf("localhost:%s", port))
			}
		}

		config.Addrs = redisAddrs

		return config
	}

	return config
}
