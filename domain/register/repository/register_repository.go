package repository

import (
	"context"
	"errors"
	"strings"

	"pokedex_backend_go/pkg/database"
	"pokedex_backend_go/pkg/model"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidEmail       = errors.New("invalid email format")
)

func NewRepository() *Repository {
	return &Repository{
		logger: zap.L().Named("register_repository"),
	}
}

type Repository struct {
	logger *zap.Logger
}

func (r *Repository) Register(ctx context.Context, email, password string) (user *model.User, err error) {
	if !isValidEmail(email) {
		r.logger.Error("Invalid email format", zap.String("email", email))
		return nil, ErrInvalidEmail
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		r.logger.Error("Failed to hash password", zap.Error(err))
		return nil, err
	}

	newUser := &model.User{
		Email:    email,
		Password: string(hashedPassword),
	}

	orm := database.Orm(ctx)

	var existingUser model.User
	result := orm.WithContext(ctx).Where("email = ?", email).First(&existingUser)
	if result.Error == nil {
		r.logger.Error("Email already exists", zap.String("email", email))
		return nil, ErrEmailAlreadyExists
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		r.logger.Error("Failed to check existing email", zap.Error(result.Error))
		return nil, result.Error
	}

	result = orm.WithContext(ctx).Create(newUser)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate") || strings.Contains(result.Error.Error(), "unique") {
			r.logger.Error("Email already exists (database constraint)", zap.String("email", email))
			return nil, ErrEmailAlreadyExists
		}
		r.logger.Error("Failed to create user", zap.Any("user", newUser), zap.Error(result.Error))
		return nil, result.Error
	}

	r.logger.Info("User created successfully", zap.String("email", email), zap.String("id", newUser.ID))

	newUser.Password = ""
	return newUser, nil
}

func isValidEmail(email string) bool {
	if email == "" {
		return false
	}

	atIndex := strings.Index(email, "@")
	if atIndex == -1 || atIndex == 0 || atIndex == len(email)-1 {
		return false
	}

	domain := email[atIndex+1:]
	if !strings.Contains(domain, ".") || strings.HasPrefix(domain, ".") || strings.HasSuffix(domain, ".") {
		return false
	}

	return true
}
