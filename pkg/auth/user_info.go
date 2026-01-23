package auth

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrNoAdminRole = status.Errorf(codes.PermissionDenied, "admin role required")

type contextKey string

var userInfoKey = contextKey("USER_INFO_KEY")

type UserInfo struct {
	ID    uuid.UUID
	roles []Role
}

func NewUserInfo(userID uuid.UUID) UserInfo {
	return UserInfo{
		ID:    userID,
		roles: getUserRoles(userID),
	}
}

var EmptyUserInfo = UserInfo{
	ID:    uuid.Nil,
	roles: []Role{},
}

func WithUserInfo(ctx context.Context, newUserInfo UserInfo) context.Context {
	return context.WithValue(ctx, &userInfoKey, newUserInfo)
}

func GetUserInfo(ctx context.Context) UserInfo {
	if userInfo, ok := ctx.Value(&userInfoKey).(UserInfo); ok {
		return userInfo
	}

	return EmptyUserInfo
}

func (ui *UserInfo) IsEmpty() bool {
	return ui.ID == uuid.Nil
}

func (ui *UserInfo) IsAdmin() bool {
	for _, userRole := range ui.roles {
		if userRole == Admin {
			return true
		}
	}

	return false
}
