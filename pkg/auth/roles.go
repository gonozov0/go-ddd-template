package auth

import (
	"github.com/google/uuid"

	"go-ddd-template/pkg/auth/tests"
)

type Role string

const Admin Role = "admin"

var userRoles = map[uuid.UUID][]Role{
	tests.AdminUserID: {Admin},
}

func getUserRoles(userID uuid.UUID) []Role {
	roles := userRoles[userID]

	if roles == nil {
		return []Role{}
	}

	return roles
}
