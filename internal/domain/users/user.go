package users

import (
	"fmt"

	"github.com/google/uuid"
)

type User struct {
	id    uuid.UUID
	name  string
	email string
}

func NewUser(name, email string) (*User, error) {
	if err := validateUsername(name); err != nil {
		return nil, err
	}
	if err := validateEmail(email); err != nil {
		return nil, err
	}

	return &User{
		id:    uuid.New(),
		name:  name,
		email: email,
	}, nil
}

func NewUserWithID(id uuid.UUID, name, email string) (*User, error) {
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

func (u *User) GetID() uuid.UUID {
	return u.id
}

func (u *User) GetName() string {
	return u.name
}

func (u *User) GetEmail() string {
	return u.email
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
