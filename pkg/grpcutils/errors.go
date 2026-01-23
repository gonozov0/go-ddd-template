package grpcutils

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"fmt"

	"go-ddd-template/pkg/testify"
)

const (
	ErrMsgInternalError   = "internal error"
	ErrMsgCanceled        = "request canceled by client"
	ErrMsgUnauthenticated = "authentication required"
)

var (
	ErrInternalError   = status.Errorf(codes.Internal, ErrMsgInternalError)
	ErrCanceled        = status.Errorf(codes.Canceled, ErrMsgCanceled)
	ErrUnauthenticated = status.Errorf(codes.Unauthenticated, ErrMsgUnauthenticated)
)

// CheckCode проверяет статус код в error и возвращает сообщение
func CheckCode(err error, code codes.Code) (string, error) {
	grpcStatus, ok := status.FromError(err)
	if !ok {
		return "", fmt.Errorf("Faied to get status code from received error (%s)", err)
	}

	if grpcStatus.Code() != code {
		return grpcStatus.Message(), fmt.Errorf(
			"Received status code: %s, expected: %s, message: %s",
			grpcStatus.Code().String(),
			code.String(),
			grpcStatus.Message(),
		)
	}

	return grpcStatus.Message(), nil
}

func CheckCodeWithSuite(s testify.Suite, err error, code codes.Code) string {
	s.T().Helper()

	grpcStatus, ok := status.FromError(err)
	s.Require().True(ok)
	s.Require().Equal(code.String(), grpcStatus.Code().String(), grpcStatus.Message())

	return grpcStatus.Message()
}

func CheckCodeAndMessageWithSuite(s testify.Suite, err error, code codes.Code, message string) {
	s.T().Helper()

	grpcStatus, ok := status.FromError(err)
	s.Require().True(ok)
	s.Require().ErrorContains(grpcStatus.Err(), message)
	s.Require().Equal(code.String(), grpcStatus.Code().String(), grpcStatus.Message())
}

func IsUnauthenticatedError(err error) bool {
	if err == nil {
		return false
	}

	if st, ok := status.FromError(err); ok {
		return st.Code() == codes.Unauthenticated
	}

	return false
}
