package orders

import (
	"context"
	"errors"
	"net/http"

	ordersDomain "go-echo-template/internal/domain/orders"

	"go-echo-template/generated/openapi"
	"go-echo-template/generated/protobuf"
	productsDomain "go-echo-template/internal/domain/products"
	usersDomain "go-echo-template/internal/domain/users"
	service "go-echo-template/internal/service/orders/create"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	// TODO: implement authentication interceptor
	order, err := ocs.CreateOrder(c.Request().Context(), c.Get("user_id").(uuid.UUID), items)
	if err != nil {
		msg := err.Error()
		if errors.Is(err, ordersDomain.ErrProductAlreadyReserved) {
			return c.JSON(http.StatusConflict, openapi.ErrorResponse{Message: &msg})
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

func (h OrderHandlers) CreateOrder(
	ctx context.Context,
	req *protobuf.CreateOrderRequest,
) (*protobuf.CreateOrderResponse, error) {
	items := make([]service.Item, 0, len(req.GetItems()))
	for _, i := range req.GetItems() {
		uid, err := uuid.Parse(i.GetId())
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid UUID: %s", i.GetId())
		}
		items = append(items, service.Item{
			ID: uid,
		})
	}

	ocs := service.NewOrderCreationService(h.orderRepo, h.userRepo, h.productRepo)
	// TODO: implement authentication interceptor
	order, err := ocs.CreateOrder(ctx, ctx.Value("user_id").(uuid.UUID), items)
	if err != nil {
		if errors.Is(err, ordersDomain.ErrProductAlreadyReserved) {
			return nil, status.Errorf(codes.Aborted, "product already reserved")
		}
		if errors.Is(err, usersDomain.ErrUserNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		if errors.Is(err, productsDomain.ErrProductNotFound) {
			return nil, status.Errorf(codes.NotFound, "product not found")
		}
		return nil, status.Errorf(codes.Internal, "internal server error")
	}
	return &protobuf.CreateOrderResponse{
		Id: order.ID().String(),
	}, nil
}
