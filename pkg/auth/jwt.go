package auth

import (
	"errors"
	"time"

	"pokedex_backend_go/pkg/model"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

const jwtSecret = "your-super-secret-jwt-key-change-this-in-production"

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type JWTService struct {
	logger *zap.Logger
}

func NewJWTService() *JWTService {
	return &JWTService{
		logger: zap.L().Named("jwt_service"),
	}
}

func (j *JWTService) GenerateToken(user *model.User) (string, error) {
	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "pokedex_backend_go",
			Subject:   user.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		j.logger.Error("Failed to sign JWT token", zap.Error(err))
		return "", err
	}

	j.logger.Info("JWT token generated successfully", zap.String("user_id", user.ID), zap.String("email", user.Email))
	return tokenString, nil
}

func (j *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		j.logger.Error("Failed to parse JWT token", zap.Error(err))
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		j.logger.Debug("JWT token validated successfully", zap.String("user_id", claims.UserID))
		return claims, nil
	}

	j.logger.Error("Invalid JWT token")
	return nil, errors.New("invalid token")
}

func (j *JWTService) RefreshToken(tokenString string) (string, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	newClaims := &Claims{
		UserID: claims.UserID,
		Email:  claims.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "pokedex_backend_go",
			Subject:   claims.UserID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	tokenString, err = token.SignedString([]byte(jwtSecret))
	if err != nil {
		j.logger.Error("Failed to refresh JWT token", zap.Error(err))
		return "", err
	}

	j.logger.Info("JWT token refreshed successfully", zap.String("user_id", claims.UserID))
	return tokenString, nil
}
