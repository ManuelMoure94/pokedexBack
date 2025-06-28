package auth

import (
	"context"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

type ContextKey string

const (
	UserContextKey ContextKey = "user"
)

type AuthMiddleware struct {
	jwtService *JWTService
	logger     *zap.Logger
}

func NewAuthMiddleware(jwtService *JWTService) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
		logger:     zap.L().Named("auth_middleware"),
	}
}

func (a *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			a.logger.Warn("Missing Authorization header")
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			a.logger.Warn("Invalid Authorization header format")
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		claims, err := a.jwtService.ValidateToken(tokenString)
		if err != nil {
			a.logger.Warn("Invalid JWT token", zap.Error(err))
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		r = r.WithContext(ctx)

		a.logger.Debug("User authenticated successfully", zap.String("user_id", claims.UserID))

		next.ServeHTTP(w, r)
	})
}

func (a *AuthMiddleware) OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenString := parts[1]

				claims, err := a.jwtService.ValidateToken(tokenString)
				if err == nil {
					ctx := context.WithValue(r.Context(), UserContextKey, claims)
					r = r.WithContext(ctx)
					a.logger.Debug("User authenticated successfully (optional)", zap.String("user_id", claims.UserID))
				} else {
					a.logger.Debug("Invalid token in optional auth", zap.Error(err))
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}

func GetUserFromContext(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(UserContextKey).(*Claims)
	return claims, ok
}
