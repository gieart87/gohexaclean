package request

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// CreateUserRequest represents the request to create a new user
type CreateUserRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

// Validate validates CreateUserRequest
func (r CreateUserRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Email,
			validation.Required.Error("email is required"),
			is.Email.Error("email must be a valid email address"),
		),
		validation.Field(&r.Name,
			validation.Required.Error("name is required"),
			validation.Length(3, 100).Error("name must be between 3 and 100 characters"),
		),
		validation.Field(&r.Password,
			validation.Required.Error("password is required"),
			validation.Length(6, 0).Error("password must be at least 6 characters"),
		),
	)
}

// UpdateUserRequest represents the request to update a user
type UpdateUserRequest struct {
	Name string `json:"name"`
}

// Validate validates UpdateUserRequest
func (r UpdateUserRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name,
			validation.Required.Error("name is required"),
			validation.Length(3, 100).Error("name must be between 3 and 100 characters"),
		),
	)
}

// LoginRequest represents the login request
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Validate validates LoginRequest
func (r LoginRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Email,
			validation.Required.Error("email is required"),
			is.Email.Error("email must be a valid email address"),
		),
		validation.Field(&r.Password,
			validation.Required.Error("password is required"),
		),
	)
}
