package suites

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	"fmt"

	"github.com/cenkalti/backoff/v4"
	"github.com/stretchr/testify/suite"

	"go-ddd-template/internal"
	pgmigrations "go-ddd-template/migrations/postgres"
	ydbmigrations "go-ddd-template/migrations/ydb"
	pgsuites "go-ddd-template/pkg/db/postgres/suites"
	redissuites "go-ddd-template/pkg/db/redis/suites"
	ydbsuites "go-ddd-template/pkg/db/ydb/suites"
	imagestorage "go-ddd-template/pkg/image_storage"
	loggerutils "go-ddd-template/pkg/logger/utils"
	"go-ddd-template/pkg/netutils"
)

type ServersInfo struct {
	PublicServer   netutils.ServerInfo
	ConsumerServer netutils.ServerInfo
	CronServer     netutils.ServerInfo
}

type runServerSuite struct {
	suite.Suite
	pgsuites.RunPostgresSuite
	redissuites.RunRedisSuite
	ydbsuites.RunYDBSuite

	wg           sync.WaitGroup
	cfg          internal.Config
	serversInfo  ServersInfo
	ImageStorage *imagestorage.Client
}

func (s *runServerSuite) SetupSuite() {
	var err error

	s.cfg, err = internal.LoadConfig()
	if err != nil {
		s.FailNow("Failed to load config", err)
	}

	s.cfg.Postgres = s.RunPostgresSuite.SetupSuite(&s.Suite, s.T(), s.cfg.Postgres, pgmigrations.FS)
	s.cfg.Redis = s.RunRedisSuite.SetupSuite(&s.Suite, s.T(), s.cfg.Redis)
	s.RunYDBSuite.SetupSuite(&s.Suite, s.T(), s.cfg.YDB, ydbmigrations.FS)

	s.ImageStorage, err = imagestorage.NewClient(context.Background(), s.cfg.ImageStorage)
	if err != nil {
		s.FailNow("failed to create image storage client: %w", err)
	}

	freePorts := make([]string, 4)
	for i := 0; i < len(freePorts); i++ {
		freePorts[i], err = netutils.GetFreePort()
		if err != nil {
			s.FailNow("Failed to get free port", err)
		}
	}

	s.serversInfo = ServersInfo{
		PublicServer:   netutils.NewServerInfo(false, "localhost", freePorts[0], freePorts[1]),
		ConsumerServer: netutils.NewServerInfo(false, "localhost", "", freePorts[2]),
		CronServer:     netutils.NewServerInfo(false, "localhost", "", freePorts[3]),
	}

	s.startCrons()
	s.startServer()
	s.startConsumers()

	s.waitToStartServers()
}

func (s *runServerSuite) waitToStartServers() {
	urls := map[string]string{
		"server":    s.serversInfo.PublicServer.GetHTTPURL(),
		"consumers": s.serversInfo.ConsumerServer.GetHTTPURL(),
		"cron":      s.serversInfo.CronServer.GetHTTPURL(),
	}

	for server, url := range urls {
		err := backoff.Retry(func() error {
			resp, err := http.Get(url + "/ping/")
			if err != nil {
				return fmt.Errorf("failed to ping: %w", err)
			}

			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("there is not OK status code")
			}

			if err := resp.Body.Close(); err != nil {
				return fmt.Errorf("failed to close response body: %w", err)
			}

			return nil
		}, backoff.WithMaxRetries(backoff.NewConstantBackOff(time.Second), 10))
		s.Require().NoErrorf(err, "failed to ping %s", server)
	}
}

func (s *runServerSuite) startServer() {
	s.cfg.Server.HTTPPort = s.serversInfo.PublicServer.HTTPPort
	s.cfg.Server.GRPCPort = s.serversInfo.PublicServer.GRPCPort
	s.cfg.Server.PprofPort = ""

	s.wg.Add(1)

	go func() {
		defer s.wg.Done()

		err := internal.RunServers(s.cfg, s.ImageStorage)
		if err != nil {
			slog.Error("Failed to run server", loggerutils.ErrAttr(fmt.Errorf("failed to run server: %w", err)))
			os.Exit(1)
		}
	}()
}

func (s *runServerSuite) startConsumers() {
	s.cfg.Consumers.HTTPPort = s.serversInfo.ConsumerServer.HTTPPort

	s.wg.Add(1)

	go func() {
		defer s.wg.Done()

		err := internal.RunConsumers(s.cfg)
		if err != nil {
			slog.Error(
				"Failed to run consumers",
				loggerutils.ErrAttr(fmt.Errorf("failed to run consumers: %w", err)),
			)
			os.Exit(1)
		}
	}()
}

func (s *runServerSuite) startCrons() {
	s.cfg.Crons.HTTPPort = s.serversInfo.CronServer.HTTPPort

	s.wg.Add(1)

	go func() {
		defer s.wg.Done()

		err := internal.RunCrons(s.cfg)
		if err != nil {
			slog.Error("Failed to run crons", loggerutils.ErrAttr(fmt.Errorf("failed to run crons: %w", err)))
			os.Exit(1)
		}
	}()
}

func (s *runServerSuite) TearDownSuite() {
	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		s.FailNow("Failed to find process", err)
	}

	err = p.Signal(os.Interrupt)
	if err != nil {
		s.FailNow("Failed to send interrupt signal", err)
	}

	s.wg.Wait()
}
