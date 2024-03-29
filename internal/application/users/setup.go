package users

import (
	"go-echo-ddd-template/internal/application/users/handlers"

	"github.com/labstack/echo/v4"
)

func Setup(e *echo.Echo, repo handlers.Repository) {
	handler := handlers.NewHandler(repo)

	userGroup := e.Group("/users")
	userGroup.POST("", handler.CreateUser)
	userGroup.GET("/:id", handler.GetUser)
}
