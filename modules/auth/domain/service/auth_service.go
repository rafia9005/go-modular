package service

import (
	"context"
	"errors"
	"go-modular-boilerplate/modules/users/domain/entity"
	"go-modular-boilerplate/modules/users/domain/repository"
)

// Errors
var (
	ErrUserNotFound     = errors.New("user not found")
	ErrEmailAlreadyUsed = errors.New("email already in use")
)

// AuthService handles user authentication
type AuthService struct {
	userRepo repository.UserRepository
}

// NewAuthService creates a new AuthService
func NewAuthService(userRepo repository.UserRepository) *AuthService {
	if userRepo == nil {
		panic("userRepo cannot be nil")
	}
	return &AuthService{
		userRepo: userRepo,
	}
}

// CreateUser creates a new user
func (s *AuthService) CreateUser(ctx context.Context, user *entity.User) error {
	if user.Email == "" || user.Password == "" {
		return errors.New("email and password cannot be empty")
	}

	existingUser, err := s.userRepo.FindByEmail(ctx, user.Email)
	if err != nil && err != repository.ERR_RECORD_NOT_FOUND {
		return err
	}
	if existingUser != nil {
		return ErrEmailAlreadyUsed
	}

	return s.userRepo.Create(ctx, user)
}
