package handlers

import (
	"errors"
	"net/http"

	productsDomain "go-echo-ddd-template/internal/domain/products"
	usersDomain "go-echo-ddd-template/internal/domain/users"
	service "go-echo-ddd-template/internal/service/orders/create"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Item struct {
	ID uuid.UUID `json:"id"`
}

type Request struct {
	Items []Item `json:"items"`
}

type Response struct {
	OrderID uuid.UUID `json:"order_id"`
}

func (h *Handler) CreateAndPay(c echo.Context) error {
	var req Request
	if err := c.Bind(&req); err != nil {
		return err
	}

	items := make([]service.Item, 0, len(req.Items))
	for _, i := range req.Items {
		items = append(items, service.Item{
			ID: i.ID,
		})
	}

	ocs := service.NewOrderCreationService(h.orderRepo, h.userRepo, h.productRepo)
	order, err := ocs.CreateOrder(c.Get("user_id").(uuid.UUID), items)
	if err != nil {
		var reservedErr *service.ProductsAlreadyReservedError
		if errors.As(err, reservedErr) {
			return c.JSON(
				http.StatusConflict,
				echo.Map{
					"message":     err.Error(),
					"product_ids": reservedErr.ProductIDs, //nolint:govet // ProductIDs is not nil
				},
			)
		}
		if errors.Is(err, usersDomain.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, echo.Map{"message": err.Error()})
		}
		if errors.Is(err, productsDomain.ErrProductNotFound) {
			return c.JSON(http.StatusNotFound, echo.Map{"message": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Internal server error"})
	}

	return c.JSON(http.StatusCreated, Response{OrderID: order.ID()})
}
