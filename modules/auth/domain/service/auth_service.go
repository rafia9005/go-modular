package service

import (
	"context"
	"errors"
	"go-modular-boilerplate/internal/pkg/utils"
	"go-modular-boilerplate/modules/users/domain/entity"
	"go-modular-boilerplate/modules/users/domain/repository"
)

// Errors
var (
	ErrUserNotFound     = errors.New("user not found")
	ErrEmailAlreadyUsed = errors.New("email already in use")
	ErrInvalidPassword  = errors.New("invalid password")
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

	// Hash the password before saving the user
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	return s.userRepo.Create(ctx, user)
}

// ProcessLogin handles user login and password verification
func (s *AuthService) ProcessLogin(ctx context.Context, email, password string) (*entity.User, error) {
	// Validate input
	if email == "" || password == "" {
		return nil, errors.New("email and password cannot be empty")
	}

	// Find user by email
	existingUser, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if err == repository.ERR_RECORD_NOT_FOUND {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// Compare the provided password with the hashed password in the database
	if !utils.CompareHashAndPassword(existingUser.Password, password) {
		return nil, ErrInvalidPassword
	}

	// Return the authenticated user
	return existingUser, nil
}
