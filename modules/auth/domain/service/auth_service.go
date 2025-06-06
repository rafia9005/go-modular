package service

import (
	"context"
	"errors"
	"go-modular/internal/pkg/jwt"
	"go-modular/internal/pkg/utils"
	"go-modular/modules/users/domain/entity"
	"go-modular/modules/users/domain/repository"
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
	jwt      jwt.JWT
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

func (s *AuthService) ChangePassword(ctx context.Context, userID uint, password string) (*entity.User, error) {
	if password == "" {
		return nil, errors.New("password cannot be empty")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	user.Password = hashedPassword

	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return nil, errors.New("failed to update password")
	}

	return user, nil
}
