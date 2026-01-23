package users

import (
	protobuf "go-ddd-template/generated/server"
	service "go-ddd-template/internal/service/users"
)

type UserHandlers struct {
	protobuf.UnimplementedUserServiceServer
	userService service.UserService
}

func SetupHandlers(userRepo service.UserRepository) UserHandlers {
	userService := service.NewUserService(userRepo)

	//nolint:exhaustivestruct
	return UserHandlers{
		userService: userService,
	}
}
