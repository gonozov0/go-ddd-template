package users

import (
	"errors"
	"fmt"
)

type Name string

var ErrInvalidName = errors.New("invalid name")

func NewName(name string) (Name, error) {
	if name == "" {
		return "", fmt.Errorf("%w: name is required", ErrInvalidName)
	}

	return Name(name), nil
}

func (n Name) String() string {
	return string(n)
}
