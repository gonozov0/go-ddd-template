package users

import (
	"errors"

	"go-ddd-template/internal/domain/shared/valueobjects"
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrInvalidUser      = errors.New("invalid user")
	ErrUserValidation   = errors.New("validation error")
	ErrUserAlreadyExist = errors.New("user already exist")
)

type User struct {
	id    valueobjects.UserID
	name  Name
	email valueobjects.Email
}

func NewUser(id valueobjects.UserID, name Name, email valueobjects.Email) *User {
	return &User{
		id:    id,
		name:  name,
		email: email,
	}
}

func (u *User) Update(name Name, email valueobjects.Email) {
	u.name = name
	u.email = email
}

func (u *User) GetID() valueobjects.UserID {
	return u.id
}

func (u *User) GetName() Name {
	return u.name
}

func (u *User) GetEmail() valueobjects.Email {
	return u.email
}
