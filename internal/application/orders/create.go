package orders

import (
	"errors"
	"net/http"

	"go-echo-ddd-template/generated/openapi"
	productsDomain "go-echo-ddd-template/internal/domain/products"
	usersDomain "go-echo-ddd-template/internal/domain/users"
	service "go-echo-ddd-template/internal/service/orders/create"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (h OrderHandlers) PostOrders(c echo.Context) error {
	var req openapi.CreateOrderRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	items := make([]service.Item, 0, len(req.Items))
	for _, i := range req.Items {
		items = append(items, service.Item{
			ID: *i.Id,
		})
	}

	ocs := service.NewOrderCreationService(h.orderRepo, h.userRepo, h.productRepo)
	order, err := ocs.CreateOrder(c.Get("user_id").(uuid.UUID), items)
	if err != nil {
		msg := err.Error()
		var reservedErr *service.ProductsAlreadyReservedError
		if errors.As(err, reservedErr) {
			return c.JSON(
				http.StatusConflict,
				openapi.ConflictOrderResponse{
					Message:    &msg,
					ProductIds: &reservedErr.ProductIDs, //nolint:govet // ProductIDs is not nil
				},
			)
		}
		if errors.Is(err, usersDomain.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, openapi.ErrorResponse{Message: &msg})
		}
		if errors.Is(err, productsDomain.ErrProductNotFound) {
			return c.JSON(http.StatusNotFound, openapi.ErrorResponse{Message: &msg})
		}
		return c.JSON(http.StatusInternalServerError, openapi.ErrorResponse{Message: &msg})
	}
	id := order.ID()
	return c.JSON(http.StatusCreated, openapi.CreateOrderResponse{Id: &id})
}
