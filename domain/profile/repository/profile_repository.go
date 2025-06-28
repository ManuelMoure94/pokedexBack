package repository

import (
	"context"
	"errors"

	"pokedex_backend_go/pkg/database"
	"pokedex_backend_go/pkg/model"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

var ErrUserNotFound = errors.New("user not found")

func NewRepository() *Repository {
	return &Repository{
		logger: zap.L().Named("profile_repository"),
	}
}

type Repository struct {
	logger *zap.Logger
}

func (r *Repository) GetUserByID(ctx context.Context, userID string) (user *model.User, err error) {
	orm := database.Orm(ctx)

	var foundUser model.User
	result := orm.WithContext(ctx).Where("id = ?", userID).First(&foundUser)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			r.logger.Error("User not found", zap.String("user_id", userID))
			return nil, ErrUserNotFound
		}
		r.logger.Error("Failed to find user", zap.String("user_id", userID), zap.Error(result.Error))
		return nil, result.Error
	}

	r.logger.Info("User found successfully", zap.String("user_id", userID), zap.String("email", foundUser.Email))

	foundUser.Password = ""
	return &foundUser, nil
}

func (r *Repository) UpdateUser(ctx context.Context, userID string, updates map[string]interface{}) (user *model.User, err error) {
	orm := database.Orm(ctx)

	var foundUser model.User
	result := orm.WithContext(ctx).Where("id = ?", userID).First(&foundUser)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			r.logger.Error("User not found for update", zap.String("user_id", userID))
			return nil, ErrUserNotFound
		}
		r.logger.Error("Failed to find user for update", zap.String("user_id", userID), zap.Error(result.Error))
		return nil, result.Error
	}

	if username, exists := updates["username"]; exists {
		if usernameStr, ok := username.(string); ok && usernameStr != "" {
			var existingUser model.User
			result := orm.WithContext(ctx).Where("username = ? AND id != ?", usernameStr, userID).First(&existingUser)
			if result.Error == nil {
				r.logger.Error("Username already exists", zap.String("username", usernameStr))
				return nil, errors.New("username already exists")
			} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
				r.logger.Error("Failed to check existing username", zap.Error(result.Error))
				return nil, result.Error
			}
		}
	}

	result = orm.WithContext(ctx).Model(&foundUser).Updates(updates)
	if result.Error != nil {
		r.logger.Error("Failed to update user", zap.String("user_id", userID), zap.Error(result.Error))
		return nil, result.Error
	}

	r.logger.Info("User updated successfully", zap.String("user_id", userID), zap.Any("updates", updates))

	foundUser.Password = ""
	return &foundUser, nil
}
