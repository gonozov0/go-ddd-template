package handlers

import (
	"errors"
	"net/http"

	productsDomain "go-echo-ddd-template/internal/domain/products"
	usersDomain "go-echo-ddd-template/internal/domain/users"
	service "go-echo-ddd-template/internal/service/orders/create"
	"go-echo-ddd-template/pkg/responses"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Item struct {
	ID uuid.UUID `json:"id"`
}

type CreateRequest struct {
	Items []Item `json:"items"`
}

type CreateResponse struct {
	ID uuid.UUID `json:"id"`
}

type ConflictResponse struct {
	Message    string      `json:"message"`
	ProductIDs []uuid.UUID `json:"product_ids"`
}

// CreateAndPay creates and processes payment for an order.
//
//	@Summary		Create and pay for an order
//	@Description	Creates an order with the items provided and processes payment
//	@Tags			orders
//	@Accept			json
//	@Produce		json
//	@Param			request	body		CreateRequest			true	"Order creation request"
//	@Success		201		{object}	CreateResponse			"Order successfully created"
//	@Failure		400		{object}	responses.ErrorResponse	"Invalid request data"
//	@Failure		404		{object}	responses.ErrorResponse	"User or product not found"
//	@Failure		409		{object}	ConflictResponse		"Products already reserved"
//	@Failure		500		{object}	responses.ErrorResponse	"Internal server error"
//	@Router			/orders [post]
func (h *Handler) CreateAndPay(c echo.Context) error {
	var req CreateRequest
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
				ConflictResponse{
					Message:    err.Error(),
					ProductIDs: reservedErr.ProductIDs, //nolint:govet // ProductIDs is not nil
				},
			)
		}
		if errors.Is(err, usersDomain.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, responses.ErrorResponse{Message: err.Error()})
		}
		if errors.Is(err, productsDomain.ErrProductNotFound) {
			return c.JSON(http.StatusNotFound, responses.ErrorResponse{Message: err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Message: "Internal server error"})
	}

	return c.JSON(http.StatusCreated, CreateResponse{ID: order.ID()})
}
