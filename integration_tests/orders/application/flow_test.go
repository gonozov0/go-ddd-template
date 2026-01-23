package application_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"google.golang.org/grpc/codes"

	pb "go-ddd-template/generated/server"
	userunithelpers "go-ddd-template/internal/application/users/shared/helpers"
	"go-ddd-template/internal/domain/shared/valueobjects"
	"go-ddd-template/pkg/grpcutils"
)

func (s *OrdersSuite) TestGRPC() {
	var orderID valueobjects.OrderID

	s.userHelper.CreateUser(userunithelpers.UserWithID(s.UserID))

	productIDs := s.productHelper.CreateRandomProducts(3)

	s.Run("Create Order", func() {
		req := &pb.CreateOrderRequest{
			ProductIds: productIDs.Strings(),
		}

		resp, err := s.grpcClient.CreateOrder(s.UserCtx, req)
		grpcutils.CheckCodeWithSuite(s, err, codes.OK)
		s.Require().NotNil(resp)

		orderID, err = valueobjects.NewOrderIDFromString(resp.GetId())
		s.Require().NoError(err)
		s.Require().NotEmpty(orderID)

		deliveryID := s.deliveryHelper.WaitForDeliveryCreation(orderID)
		s.Require().NotEmpty(deliveryID)
	})

	s.T().Cleanup(func() {
		s.orderHelper.DeleteOrderInCleanup(orderID)
	})
}

func (s *OrdersSuite) TestHTTP() {
	var orderID valueobjects.OrderID

	s.userHelper.CreateUser(userunithelpers.UserWithID(s.UserID))

	productIDs := s.productHelper.CreateRandomProducts(3)

	s.Run("Create Order", func() {
		req := &pb.CreateOrderRequest{
			ProductIds: productIDs.Strings(),
		}

		body, err := grpcutils.MarshalJSON(req)
		s.Require().NoError(err)

		httpReq, err := http.NewRequest(
			http.MethodPost,
			fmt.Sprintf("%s/orders", s.ServerURL),
			bytes.NewBuffer(body),
		)
		s.Require().NoError(err)

		httpReq.Header = s.UserHeaders

		httpResp, err := http.DefaultClient.Do(httpReq)
		s.Require().NoError(err)

		respBytes, err := io.ReadAll(httpResp.Body)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, httpResp.StatusCode, string(respBytes))

		var resp pb.CreateOrderResponse

		err = grpcutils.UnmarshalJSON(respBytes, &resp)
		s.Require().NoError(err)
		s.Require().NoError(httpResp.Body.Close())

		orderID, err = valueobjects.NewOrderIDFromString(resp.GetId())
		s.Require().NoError(err)
		s.Require().NotEmpty(orderID)

		deliveryID := s.deliveryHelper.WaitForDeliveryCreation(orderID)
		s.Require().NotEmpty(deliveryID)
	})

	s.T().Cleanup(func() {
		s.orderHelper.DeleteOrderInCleanup(orderID)
	})
}
