package auth

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"go-ddd-template/pkg/grpcutils"
)

var (
	ErrMissingUserMetadata = status.Errorf(codes.Unauthenticated, "missing user metadata")
	ErrInvalidUserID       = status.Errorf(codes.InvalidArgument, "invalid user id")
)

type authMiddleware struct {
	cfg Config
}

// NewAuthMiddleware creates a new auth middleware to check tvm tickets.
func NewAuthMiddleware(cfg Config) (grpc.UnaryServerInterceptor, error) {
	mw := &authMiddleware{
		cfg: cfg,
	}

	return mw.UnaryInterceptor, nil
}

//nolint:cyclop
func (mw *authMiddleware) auth(_ *grpc.UnaryServerInfo, md metadata.MD) (UserInfo, error) {
	rawUserID, err := grpcutils.GetSingleHeader(md, mw.cfg.UserHeader)
	if err != nil || rawUserID == "" {
		return EmptyUserInfo, nil
	}

	userID, err := uuid.Parse(rawUserID)
	if err != nil {
		return EmptyUserInfo, ErrInvalidUserID
	}

	return NewUserInfo(userID), nil
}

func (mw *authMiddleware) UnaryInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, ErrMissingUserMetadata
	}

	if IsAuthRequired(info.FullMethod) {
		userInfo, err := mw.auth(info, md)
		if err != nil {
			return nil, err
		}

		ctx = WithUserInfo(ctx, userInfo)
	}

	return handler(ctx, req)
}
