package users

import (
	"errors"
)

var ErrUserNotFound = errors.New("user not found")
var ErrInvalidUser = errors.New("invalid user")
var ErrUserValidation = errors.New("validation error")
