package orders

import (
	"go-echo-ddd-template/internal/application/orders/handlers"

	"github.com/labstack/echo/v4"
)

func Setup(
	e *echo.Echo,
	orderRepo handlers.OrderRepository,
	userRepo handlers.UserRepository,
	productRepo handlers.ProductRepository,
) {
	handler := handlers.NewHandler(orderRepo, userRepo, productRepo)

	orderGroup := e.Group("/orders")
	orderGroup.POST("/create-and-pay", handler.CreateAndPay)
}
