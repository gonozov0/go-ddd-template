package application_test

import (
	"context"
	"encoding/json"
	"io"

	"go-echo-template/generated/openapi"
	"go-echo-template/generated/protobuf"

	openapi_types "github.com/oapi-codegen/runtime/types"
)

func (s *UsersSuite) TestCreateAndGetUser() {
	s.Run("HTTP", func() {
		var userID openapi_types.UUID
		user := openapi.CreateUserRequest{
			Name:  "John Doe",
			Email: "johndoe@example.com",
		}

		s.Run("Create", func() {
			createResp, err := s.httpClient.PostUsers(context.Background(), user)
			s.Require().NoError(err, "Failed to create user via HTTP")
			s.Require().Equal(201, createResp.StatusCode, "Failed to create user via HTTP")

			body, err := io.ReadAll(createResp.Body)
			s.Require().NoError(err, "Failed to read response body")
			defer createResp.Body.Close()
			var resp openapi.CreateUserResponse
			err = json.Unmarshal(body, &resp)
			s.Require().NoError(err, "Failed to unmarshal response body")
			s.Require().NotEmpty(resp.Id, "Received empty user ID")

			userID = *resp.Id
		})

		s.Run("Get", func() {
			getResp, err := s.httpClient.GetUsersId(context.Background(), userID)
			s.Require().NoError(err, "Failed to get user via HTTP")
			s.Require().Equal(200, getResp.StatusCode, "Failed to get user via HTTP")

			body, err := io.ReadAll(getResp.Body)
			s.Require().NoError(err, "Failed to read response body")
			defer getResp.Body.Close()
			var resp openapi.GetUserResponse
			err = json.Unmarshal(body, &resp)
			s.Require().NoError(err, "Failed to unmarshal response body")
			s.Require().Equal(user.Name, *resp.Name, "User names do not match")
			s.Require().Equal(user.Email, *resp.Email, "User emails do not match")
			s.Require().Equal(userID, *resp.Id, "User IDs do not match")
		})
	})

	s.Run("GRPC", func() {
		createReq := &protobuf.CreateUserRequest{
			Name:  "John Doe",
			Email: "johndoe@example.com",
		}

		createResp, err := s.grpcClient.CreateUser(context.Background(), createReq)

		s.Require().NoError(err)
		s.Require().NotNil(createResp)
		s.Require().NotEmpty(createResp.GetId())

		getReq := &protobuf.GetUserRequest{
			Id: createResp.GetId(),
		}
		getResp, err := s.grpcClient.GetUser(context.Background(), getReq)

		s.Require().NoError(err)
		s.Require().NotNil(getResp)
		s.Require().Equal(createReq.GetName(), getResp.GetName())
		s.Require().Equal(createReq.GetEmail(), getResp.GetEmail())
		s.Require().Equal(createResp.GetId(), getResp.GetId())
	})
}
