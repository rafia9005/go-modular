package request

// LoginRequest represents a request to login a user
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// CreateUserRequest represents a request to create a user
type CreateUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// UpdateUserRequest represents a request to update a user
type UpdateUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"omitempty,min=6"`
}

type ChnagePasswordRequest struct {
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required, min=6"`
}
