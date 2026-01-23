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

func (s *UsersSuite) TestGRPC() {
	var (
		user = userunithelpers.GenerateUser(&s.BaseSuite)
	)

	s.Run("Create User", func() {
		resp, err := s.grpcClient.CreateUser(s.AdminCtx, userunithelpers.ToCreateUserRequest(user))
		grpcutils.CheckCodeWithSuite(s, err, codes.OK)
		s.Require().NotNil(resp)

		s.Require().Equal(user.GetID().String(), resp.GetId())

		userID, err := valueobjects.NewUserIDFromString(resp.GetId())
		s.Require().NoError(err)

		user, err = userunithelpers.UpdateUserWithID(userID, user)
		s.Require().NoError(err)
	})

	s.T().Cleanup(func() {
		s.userHelper.DeleteUserInCleanup(user.GetID())
	})

	s.Run("Get User", func() {
		s.Require().NotEmpty(user.GetID())

		getReq := &pb.GetUserRequest{Id: user.GetID().String()}

		resp, err := s.grpcClient.GetUser(s.AdminCtx, getReq)
		grpcutils.CheckCodeWithSuite(s, err, codes.OK)
		s.Require().NotNil(resp)

		s.Require().Equal(user.GetID().String(), resp.GetId())
		s.Require().Equal(user.GetName().String(), resp.GetName())
		s.Require().Equal(user.GetEmail().String(), resp.GetEmail())
	})
}

func (s *UsersSuite) TestHTTP() {
	var (
		user = userunithelpers.GenerateUser(s)
	)

	s.Run("Create User", func() {
		createUserBody := userunithelpers.ToCreateUserRequest(user)

		body, err := grpcutils.MarshalJSON(createUserBody)
		s.Require().NoError(err)

		httpReq, err := http.NewRequest(
			http.MethodPost,
			fmt.Sprintf("%s/users", s.ServerURL),
			bytes.NewBuffer(body),
		)
		s.Require().NoError(err)

		httpReq.Header = s.AdminHeaders

		httpResp, err := http.DefaultClient.Do(httpReq)
		s.Require().NoError(err)

		respBytes, err := io.ReadAll(httpResp.Body)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, httpResp.StatusCode, string(respBytes))

		var resp pb.CreateUserResponse

		err = grpcutils.UnmarshalJSON(respBytes, &resp)
		s.Require().NoError(err)
		s.Require().NoError(httpResp.Body.Close())

		s.Require().Equal(user.GetID().String(), resp.Id)

		userID, err := valueobjects.NewUserIDFromString(resp.Id)
		s.Require().NoError(err)

		user, err = userunithelpers.UpdateUserWithID(userID, user)
		s.Require().NoError(err)
	})

	s.T().Cleanup(func() {
		s.userHelper.DeleteUserInCleanup(user.GetID())
	})

	s.Run("Get User", func() {
		s.Require().NotEmpty(user.GetID())

		httpReq, err := http.NewRequest(
			http.MethodGet,
			fmt.Sprintf("%s/users/%s", s.ServerURL, user.GetID().String()),
			nil,
		)
		s.Require().NoError(err)

		httpReq.Header = s.AdminHeaders

		httpResp, err := http.DefaultClient.Do(httpReq)
		s.Require().NoError(err)

		respBytes, err := io.ReadAll(httpResp.Body)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, httpResp.StatusCode, string(respBytes))

		var resp pb.GetUserResponse

		err = grpcutils.UnmarshalJSON(respBytes, &resp)
		s.Require().NoError(err)
		s.Require().NoError(httpResp.Body.Close())

		s.Require().Equal(user.GetID().String(), resp.Id)
		s.Require().Equal(user.GetName().String(), resp.Name)
		s.Require().Equal(user.GetEmail().String(), resp.Email)
	})
}
