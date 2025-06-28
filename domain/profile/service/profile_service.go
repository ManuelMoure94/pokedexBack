package service

import (
	"context"
	"errors"
	"strings"

	"pokedex_backend_go/domain/profile/repository"
	"pokedex_backend_go/pkg/model"

	"go.uber.org/zap"
)

func NewService(repo *repository.Repository) *Service {
	return &Service{
		logger: zap.L().Named("profileService"),
		repo:   repo,
	}
}

type Service struct {
	logger *zap.Logger
	repo   *repository.Repository
}

func (s *Service) GetProfile(ctx context.Context, userID string) (user *model.User, err error) {
	if userID == "" {
		s.logger.Error("User ID is required")
		return nil, errors.New("user ID is required")
	}

	user, err = s.repo.GetUserByID(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get user profile", zap.String("user_id", userID), zap.Error(err))
		return nil, err
	}

	s.logger.Info("User profile retrieved successfully", zap.String("user_id", userID))
	return user, nil
}

func (s *Service) UpdateProfile(ctx context.Context, userID string, updates map[string]interface{}) (user *model.User, err error) {
	if userID == "" {
		s.logger.Error("User ID is required")
		return nil, errors.New("user ID is required")
	}

	if len(updates) == 0 {
		s.logger.Error("No updates provided")
		return nil, errors.New("no updates provided")
	}

	allowedFields := map[string]bool{
		"name":     true,
		"phone":    true,
		"username": true,
	}

	validUpdates := make(map[string]interface{})
	for field, value := range updates {
		if !allowedFields[field] {
			s.logger.Error("Field not allowed for update", zap.String("field", field))
			return nil, errors.New("field '" + field + "' is not allowed for update")
		}

		if strValue, ok := value.(string); ok {
			strValue = strings.TrimSpace(strValue)
			if strValue != "" {
				validUpdates[field] = strValue
			}
		} else {
			validUpdates[field] = value
		}
	}

	if len(validUpdates) == 0 {
		s.logger.Error("No valid updates provided")
		return nil, errors.New("no valid updates provided")
	}

	user, err = s.repo.UpdateUser(ctx, userID, validUpdates)
	if err != nil {
		s.logger.Error("Failed to update user profile", zap.String("user_id", userID), zap.Error(err))
		return nil, err
	}

	s.logger.Info("User profile updated successfully", zap.String("user_id", userID), zap.Any("updates", validUpdates))
	return user, nil
}
