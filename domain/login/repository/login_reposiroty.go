package repository

import (
	"context"
	"errors"

	"pokedex_backend_go/pkg/database"
	"pokedex_backend_go/pkg/model"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var ErrInvalidCredentials = errors.New("invalid username or password")

func NewRepository() *Repository {
	return &Repository{
		logger: zap.L().Named("login_repository"),
	}
}

type Repository struct {
	logger *zap.Logger
}

func (r *Repository) Login(ctx context.Context, email, password string) (user *model.User, err error) {
	orm := database.Orm(ctx)

	var foundUser model.User
	result := orm.WithContext(ctx).Where("email = ?", email).First(&foundUser)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			r.logger.Error("User not found", zap.String("email", email))
			return nil, ErrInvalidCredentials
		}
		r.logger.Error("Failed to find user", zap.String("email", email), zap.Error(result.Error))
		return nil, result.Error
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(password))
	if err != nil {
		r.logger.Error("Invalid password", zap.String("email", email))
		return nil, ErrInvalidCredentials
	}

	r.logger.Info("User login successful", zap.String("email", email), zap.String("id", foundUser.ID))

	foundUser.Password = ""
	return &foundUser, nil
}
