package service

import (
	"context"
	"errors"

	"pokedex_backend_go/domain/register/repository"
	"pokedex_backend_go/pkg/auth"
	"pokedex_backend_go/pkg/model"

	"go.uber.org/zap"
)

func NewService(repo *repository.Repository) *Service {
	return &Service{
		logger:     zap.L().Named("registerService"),
		repo:       repo,
		jwtService: auth.NewJWTService(),
	}
}

type Service struct {
	logger     *zap.Logger
	repo       *repository.Repository
	jwtService *auth.JWTService
}

func (s *Service) Register(ctx context.Context, email, password string) (user *model.User, err error) {
	if email == "" {
		s.logger.Error("Email is required")
		return nil, errors.New("email is required")
	}

	if password == "" {
		s.logger.Error("Password is required")
		return nil, errors.New("password is required")
	}

	if len(password) < 6 {
		s.logger.Error("Password too short", zap.Int("length", len(password)))
		return nil, errors.New("password must be at least 6 characters long")
	}

	user, err = s.repo.Register(ctx, email, password)
	if err != nil {
		s.logger.Error("Failed to register user", zap.String("email", email), zap.Error(err))
		return nil, err
	}

	s.logger.Info("User registered successfully", zap.String("email", email), zap.String("id", user.ID))
	return user, nil
}

func (s *Service) RegisterWithToken(ctx context.Context, email, password string) (user *model.User, token string, err error) {
	user, err = s.Register(ctx, email, password)
	if err != nil {
		return nil, "", err
	}

	token, err = s.jwtService.GenerateToken(user)
	if err != nil {
		s.logger.Error("Failed to generate JWT token", zap.String("email", email), zap.Error(err))
		return nil, "", err
	}

	return user, token, nil
}
