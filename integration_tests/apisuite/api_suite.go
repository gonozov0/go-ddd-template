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

const waitServerTimeout = time.Second * 2

type APITestSuite struct {
	suite.Suite
	ServerURL string
	Conn      *grpc.ClientConn
}

func (suite *APITestSuite) SetupSuite() {
	suite.ServerURL = "http://localhost:8080/"

	go func() {
		err := internal.Run()
		if err != nil {
			slog.Error("Error while running the server", "err", err)
			os.Exit(1)
		}
	}()

	time.Sleep(waitServerTimeout)

	var err error
	suite.Conn, err = grpc.NewClient(
		suite.ServerURL,
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
		slog.Error("Failed to find process", "err", err)
		os.Exit(1)
	}

	err = p.Signal(os.Interrupt)
	if err != nil {
		slog.Error("Failed to send interrupt signal", "err", err)
		os.Exit(1)
	}

	time.Sleep(waitServerTimeout)
}
