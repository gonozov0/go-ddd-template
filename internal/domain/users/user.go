package users

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrInvalidUser    = errors.New("invalid user")
	ErrUserValidation = errors.New("validation error")
)

type User struct {
	id    uuid.UUID
	name  string
	email string
}

func NewUser(id uuid.UUID, name, email string) (*User, error) {
	if err := validateUsername(name); err != nil {
		return nil, err
	}
	if err := validateEmail(email); err != nil {
		return nil, err
	}

	return &User{
		id:    id,
		name:  name,
		email: email,
	}, nil
}

func CreateUser(name, email string) (*User, error) {
	return NewUser(uuid.New(), name, email)
}

func (u *User) ID() uuid.UUID {
	return u.id
}

func (u *User) Name() string {
	return u.name
}

func (u *User) Email() string {
	return u.email
}

func (u *User) SendToEmail(_ string) error {
	return errors.New("not implemented")
}

func validateUsername(username string) error {
	if username == "" {
		return fmt.Errorf("%w: name is required", ErrUserValidation)
	}
	return nil
}

func validateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("%w: email is required", ErrUserValidation)
	}
	return nil
}
