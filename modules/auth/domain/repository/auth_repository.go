package repository

import (
	"context"
	"errors"
	"go-modular-boilerplate/internal/pkg/database"
	"go-modular-boilerplate/modules/users/domain/entity"
)

var (
	ERR_RECORD_NOT_FOUND = errors.New("record not found")
)

type AuthRepositoryImpl struct {}

func (r AuthRepositoryImpl) RegisterUser(ctx context.Context, user *entity.User) error {
	return database.DB.WithContext(ctx).Create(user).Error
}

func (r AuthRepositoryImpl) FindMyEmail(ctx context.Context, email string) (*entity.User, error) {
	return nil
}
