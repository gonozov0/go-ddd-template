package apisuite

import (
	"log/slog"
	"os"
	"time"

	"go-echo-ddd-template/internal"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	waitServerTimeout = 2 * time.Second
)

type APITestSuite struct {
	suite.Suite
	HTTPServerURL string
	GRPCServerURL string
	Conn          *grpc.ClientConn
}

func (suite *APITestSuite) SetupSuite() {
	suite.GRPCServerURL = "localhost:8080"
	suite.HTTPServerURL = "http://" + suite.GRPCServerURL

	go func() {
		err := internal.Run()
		if err != nil {
			suite.Fail("Failed to run server", err)
		}
	}()

	time.Sleep(waitServerTimeout)

	var err error
	suite.Conn, err = grpc.NewClient(
		suite.GRPCServerURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		slog.Error("Failed to dial server", "err", err)
		os.Exit(1)
	}
}

func (suite *APITestSuite) TearDownSuite() {
	if suite.Conn != nil {
		suite.Conn.Close()
	}

	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		suite.Fail("Failed to find process", err)
	}

	err = p.Signal(os.Interrupt)
	if err != nil {
		suite.Fail("Failed to send interrupt signal", err)
	}

	time.Sleep(waitServerTimeout)
}
