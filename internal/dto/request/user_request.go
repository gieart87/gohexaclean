package request

// CreateUserRequest represents the request to create a new user
type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required,min=3,max=100"`
	Password string `json:"password" validate:"required,min=6"`
}

// UpdateUserRequest represents the request to update a user
type UpdateUserRequest struct {
	Name string `json:"name" validate:"required,min=3,max=100"`
}

// LoginRequest represents the login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
