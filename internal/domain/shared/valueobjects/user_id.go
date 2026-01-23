package valueobjects

import (
	"errors"

	"github.com/google/uuid"
)

type UserID uuid.UUID

var EmptyUserID = UserID(uuid.Nil)

var ErrInvalidUserID = errors.New("invalid user id")

func NewUserID(id uuid.UUID) (UserID, error) {
	if id == uuid.Nil {
		return EmptyUserID, ErrInvalidUserID
	}

	return UserID(id), nil
}

func NewUserIDFromString(id string) (UserID, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return EmptyUserID, ErrInvalidUserID
	}

	return NewUserID(uid)
}

func (id UserID) UUID() uuid.UUID {
	return uuid.UUID(id)
}

func (id UserID) String() string {
	return id.UUID().String()
}

func (id UserID) IsEmpty() bool {
	return id == EmptyUserID
}
