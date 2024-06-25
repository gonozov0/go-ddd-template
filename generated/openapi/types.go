// Package openapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.3.0 DO NOT EDIT.
package openapi

import (
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// ConflictOrderResponse defines model for ConflictOrderResponse.
type ConflictOrderResponse struct {
	Message    *string               `json:"message,omitempty"`
	ProductIds *[]openapi_types.UUID `json:"product_ids,omitempty"`
}

// CreateOrderRequest defines model for CreateOrderRequest.
type CreateOrderRequest struct {
	Items []OrderItem `json:"items"`
}

// CreateOrderResponse defines model for CreateOrderResponse.
type CreateOrderResponse struct {
	Id *openapi_types.UUID `json:"id,omitempty"`
}

// CreateUserRequest defines model for CreateUserRequest.
type CreateUserRequest struct {
	Email openapi_types.Email `json:"email"`
	Name  string              `json:"name"`
}

// CreateUserResponse defines model for CreateUserResponse.
type CreateUserResponse struct {
	Id *openapi_types.UUID `json:"id,omitempty"`
}

// ErrorResponse defines model for ErrorResponse.
type ErrorResponse struct {
	Message *string `json:"message,omitempty"`
}

// GetUserResponse defines model for GetUserResponse.
type GetUserResponse struct {
	Email *openapi_types.Email `json:"email,omitempty"`
	Id    *openapi_types.UUID  `json:"id,omitempty"`
	Name  *string              `json:"name,omitempty"`
}

// OrderItem defines model for OrderItem.
type OrderItem struct {
	Id *openapi_types.UUID `json:"id,omitempty"`
}

// PostOrdersJSONRequestBody defines body for PostOrders for application/json ContentType.
type PostOrdersJSONRequestBody = CreateOrderRequest

// PostUsersJSONRequestBody defines body for PostUsers for application/json ContentType.
type PostUsersJSONRequestBody = CreateUserRequest