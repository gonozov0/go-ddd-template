package suites

import (
	"context"
	"net/http"
	"os"

	"go-ddd-template/internal/domain/shared/valueobjects"

	"go-ddd-template/pkg/auth"
	"go-ddd-template/pkg/auth/tests"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type BaseSuite struct {
	runServerSuite

	ServerURL  string
	ServerConn *grpc.ClientConn

	UserID       valueobjects.UserID
	UserCtx      context.Context
	UserHeaders  http.Header
	AdminCtx     context.Context
	AdminHeaders http.Header
}

func (s *BaseSuite) SetupSuite() {
	var err error

	setEnvsToRunCrons()
	s.runServerSuite.SetupSuite()

	s.ServerConn, err = s.serversInfo.PublicServer.MakeGRPCConn()
	if err != nil {
		s.FailNow("Failed to dial public server", err)
	}

	s.ServerURL = s.serversInfo.PublicServer.GetHTTPURL()

	s.UserID = valueobjects.UserID(tests.UserID)
	s.UserCtx = metadata.NewOutgoingContext(context.Background(), metadata.New(map[string]string{
		auth.DefaultUserHeader: tests.UserID.String(),
	}))
	s.UserHeaders = http.Header{
		auth.DefaultUserHeader: []string{tests.UserID.String()},
	}
	s.AdminCtx = metadata.NewOutgoingContext(context.Background(), metadata.New(map[string]string{
		auth.DefaultUserHeader: tests.AdminUserID.String(),
	}))
	s.AdminHeaders = http.Header{
		auth.DefaultUserHeader: []string{tests.AdminUserID.String()},
	}
}

func (s *BaseSuite) TearDownSuite() {
	if s.ServerConn != nil {
		if err := s.ServerConn.Close(); err != nil {
			s.FailNow("Failed to close connection", err)
		}
	}

	s.runServerSuite.TearDownSuite()
}

func setEnvsToRunCrons() {
	must(os.Setenv("HANDLE_CREATED_ORDERS_INTERVAL", "1s"))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
