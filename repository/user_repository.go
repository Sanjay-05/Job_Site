package repository

import (
	"context"
	"errors"
	"jobsite/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	FindUserByName(ctx context.Context, userName string) (*models.User, error)
	ResetPassword(ctx context.Context, userName, newPassword string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) FindUserByName(ctx context.Context, userName string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("user_name = ?", userName).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) ResetPassword(ctx context.Context, userName, newPassword string) error {
	result := r.db.Model(&models.User{}).Where("user_name = ?", userName).Update("password", newPassword)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("reset password failed")
	}
	return nil
}
