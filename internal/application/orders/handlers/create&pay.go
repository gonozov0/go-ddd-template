package handlers

import (
	"log/slog"
	"net/http"

	"go-echo-ddd-template/internal/service/orders"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type item struct {
	ID       uuid.UUID `json:"id"`
	Quantity int       `json:"quantity"`
}

type request struct {
	Items []item `json:"items"`
}

type response struct {
	OrderID uuid.UUID `json:"order_id"`
}

func (h *Handler) CreateAndPay(c echo.Context) error {
	var req request
	if err := c.Bind(&req); err != nil {
		return err
	}

	items := make([]orders.Item, 0, len(req.Items))
	for _, i := range req.Items {
		items = append(items, orders.Item{
			ID:       i.ID,
			Quantity: i.Quantity,
		})
	}

	os := orders.NewOrderService(h.orderRepo, h.userRepo, h.productRepo)
	order, err := os.CreateAndPay(c.Get("user_id").(uuid.UUID), items)
	if err != nil {
		slog.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Internal server error"})
	}

	return c.JSON(http.StatusCreated, response{OrderID: order.ID()})
}
