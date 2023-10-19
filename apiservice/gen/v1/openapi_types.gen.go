// Package gen provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.15.0 DO NOT EDIT.
package gen

import (
	"time"
)

// CreateUser defines model for CreateUser.
type CreateUser struct {
	Name string `json:"name"`
}

// ErrorResponse defines model for ErrorResponse.
type ErrorResponse struct {
	Message string `json:"message"`
}

// User defines model for User.
type User struct {
	CreatedAt time.Time `json:"created_at"`
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
}

// Voucher defines model for Voucher.
type Voucher struct {
	CreatedAt time.Time `json:"created_at"`
	Id        int64     `json:"id"`
	UserId    int64     `json:"user_id"`
	UserName  string    `json:"user_name"`
	Value     string    `json:"value"`
}

// SearchVoucherParams defines parameters for SearchVoucher.
type SearchVoucherParams struct {
	// UserId User ID
	UserId int64 `form:"userId" json:"userId"`
}

// CreateUserJSONRequestBody defines body for CreateUser for application/json ContentType.
type CreateUserJSONRequestBody = CreateUser
