package valueobjects

import (
	"regexp"

	"errors"
	"fmt"
)

type Email string

var (
	ErrInvalidEmail = errors.New("invalid email")
	emailRegexp     = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
)

func NewEmail(email string) (Email, error) {
	if email == "" {
		return "", fmt.Errorf("%w: email is required", ErrInvalidEmail)
	}

	if !emailRegexp.MatchString(email) {
		return "", fmt.Errorf("%w: email format is invalid", ErrInvalidEmail)
	}

	return Email(email), nil
}

func (e Email) String() string {
	return string(e)
}
