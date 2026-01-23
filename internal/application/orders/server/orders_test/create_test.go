package orders_test

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"

	pb "go-ddd-template/generated/server"
	orderunithelpers "go-ddd-template/internal/application/orders/shared/helpers"
	userunithelpers "go-ddd-template/internal/application/users/shared/helpers"
	"go-ddd-template/pkg/auth"
	"go-ddd-template/pkg/grpcutils"
)

func (s *OrdersSuite) TestCreateOrder() {
	order := s.prepareOrder()
	req := orderunithelpers.ToCreateOrderRequest(*order)

	resp, err := s.GRPCHandlers.CreateOrder(s.UserCtx, req)
	grpcutils.CheckCodeWithSuite(s, err, codes.OK)
	s.Require().NotEqual("", resp.GetId())
	s.Require().NotEqual(uuid.Nil, resp.GetId())
}

func (s *OrdersSuite) TestCreateOrderWithEmptyUserUID() {
	req := &pb.CreateOrderRequest{
		ProductIds: []string{"invalid-uuid"},
	}

	_, err := s.GRPCHandlers.CreateOrder(context.Background(), req)
	message := grpcutils.CheckCodeWithSuite(s, err, codes.Unauthenticated)

	s.Require().Contains(auth.ErrMissingUserMetadata.Error(), message)
}

func (s *OrdersSuite) TestCreateOrderWithInvalidProductIDs() {
	userunithelpers.CreateUser(s, s.UsersRepo, userunithelpers.UserWithID(s.UserID))

	req := &pb.CreateOrderRequest{
		ProductIds: []string{"invalid-uuid"},
	}
	_, err := s.GRPCHandlers.CreateOrder(s.UserCtx, req)
	message := grpcutils.CheckCodeWithSuite(s, err, codes.InvalidArgument)

	s.Require().Contains(message, "invalid product id")
}

func (s *OrdersSuite) TestCreateOrderWithUserNotFound() {
	order := s.prepareOrder()
	req := orderunithelpers.ToCreateOrderRequest(*order)

	ctx := auth.WithUserInfo(context.Background(), auth.NewUserInfo(uuid.New()))
	_, err := s.GRPCHandlers.CreateOrder(ctx, req)
	message := grpcutils.CheckCodeWithSuite(s, err, codes.NotFound)

	s.Require().Contains(message, "user not found")
}

func (s *OrdersSuite) TestCreateOrderWithProductNotFound() {
	userunithelpers.CreateUser(s, s.UsersRepo, userunithelpers.UserWithID(s.UserID))

	req := &pb.CreateOrderRequest{
		ProductIds: []string{uuid.New().String()},
	}
	_, err := s.GRPCHandlers.CreateOrder(s.UserCtx, req)
	message := grpcutils.CheckCodeWithSuite(s, err, codes.NotFound)

	s.Require().Contains(message, "product not found")
}
