package helpers

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"

	"go-ddd-template/internal/domain/shared/valueobjects"
	domain "go-ddd-template/internal/domain/users"
	"go-ddd-template/pkg/testify"
)

// Методы для генерации User
type (
	userToCreate struct {
		id    valueobjects.UserID
		name  domain.Name
		email valueobjects.Email
	}

	GenerateUserOption func(*userToCreate) error
)

func UserWithID(id valueobjects.UserID) GenerateUserOption {
	return func(user *userToCreate) error {
		user.id = id
		return nil
	}
}

func UserWithName(name domain.Name) GenerateUserOption {
	return func(user *userToCreate) error {
		user.name = name
		return nil
	}
}

func UserWithEmail(email valueobjects.Email) GenerateUserOption {
	return func(user *userToCreate) error {
		user.email = email
		return nil
	}
}

// GenerateUser генерирует User с помощью переданных опций
// Если опции не переданы, генерируется User со случайными данными:
//   - id		- случайный число от 1 до 1_000_000
//   - name		- случайное имя
//   - email	- случайный email
func GenerateUser(s testify.Suite, opts ...GenerateUserOption) *domain.User {
	user := &userToCreate{
		id:    valueobjects.UserID(uuid.New()),
		name:  domain.Name(gofakeit.Name()),
		email: valueobjects.Email(gofakeit.Email()),
	}

	for _, opt := range opts {
		s.Require().NoError(opt(user))
	}

	return domain.NewUser(user.id, user.name, user.email)
}
