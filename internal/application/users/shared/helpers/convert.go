package helpers

import (
	pb "go-ddd-template/generated/server"
	"go-ddd-template/internal/domain/shared/valueobjects"
	domain "go-ddd-template/internal/domain/users"
)

func ToCreateUserRequest(user *domain.User) *pb.CreateUserRequest {
	return &pb.CreateUserRequest{
		Id:    user.GetID().String(),
		Name:  user.GetName().String(),
		Email: user.GetEmail().String(),
	}
}

func UpdateUserWithID(userID valueobjects.UserID, user *domain.User) (*domain.User, error) {
	return domain.NewUser(
		userID,
		user.GetName(),
		user.GetEmail(),
	), nil
}
