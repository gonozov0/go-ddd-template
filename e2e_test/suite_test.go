package e2e_test

import (
	"log/slog"
	"os"
	"testing"
	"time"

	"go-echo-ddd-template/internal"

	"github.com/stretchr/testify/suite"
)

type APITestSuite struct {
	suite.Suite
	serverURL string
}

func (suite *APITestSuite) SetupSuite() {
	suite.serverURL = "http://localhost:8080/"

	go func() {
		err := internal.Run()
		if err != nil {
			slog.Error("Error while running the server", "err", err)
			os.Exit(1)
		}
	}()

	// Wait for the server to start
	time.Sleep(time.Second * 2)
}

func (suite *APITestSuite) TearDownSuite() {
	p, err := os.FindProcess(os.Getpid())
	suite.Require().NoError(err)

	err = p.Signal(os.Interrupt)
	suite.Require().NoError(err)

	// Wait for the server to stop
	time.Sleep(time.Second * 2)
}

func TestAPISuite(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}
