package service

import (
	"context"
	"errors"

	repository "pokedex_backend_go/domain/login/repository"
	"pokedex_backend_go/pkg/auth"
	"pokedex_backend_go/pkg/model"

	"go.uber.org/zap"
)

func NewService(repo *repository.Repository) *Service {
	return &Service{
		logger:     zap.L().Named("loginService"),
		repo:       repo,
		jwtService: auth.NewJWTService(),
	}
}

type Service struct {
	logger     *zap.Logger
	repo       *repository.Repository
	jwtService *auth.JWTService
}

func (s *Service) Login(ctx context.Context, email, password string) (user *model.User, err error) {
	if email == "" {
		s.logger.Error("Email is required")
		return nil, errors.New("email is required")
	}

	if password == "" {
		s.logger.Error("Password is required")
		return nil, errors.New("password is required")
	}

	userData, err := s.repo.Login(ctx, email, password)
	if err != nil {
		s.logger.Error("Failed to login", zap.String("email", email), zap.Error(err))
		return nil, err
	}

	s.logger.Info("User login successful", zap.String("email", email), zap.String("id", userData.ID))
	return userData, nil
}

func (s *Service) LoginWithToken(ctx context.Context, email, password string) (user *model.User, token string, err error) {
	user, err = s.Login(ctx, email, password)
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
