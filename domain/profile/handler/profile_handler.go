package handler

import (
	"encoding/json"
	"net/http"

	"pokedex_backend_go/domain/profile/service"
	"pokedex_backend_go/pkg/auth"
	"pokedex_backend_go/pkg/model"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type ProfileHandler struct {
	service *service.Service
	logger  *zap.Logger
}

func NewHandler(service *service.Service) *ProfileHandler {
	return &ProfileHandler{
		service: service,
		logger:  zap.L().Named("profile_handler"),
	}
}

func Handler(service *service.Service, authMiddleware *auth.AuthMiddleware) func(chi.Router) {
	return func(r chi.Router) {
		logger := zap.L().Named("profile_handler_registration")
		logger.Info("Registering profile handler at /api/v1/profile")

		handler := NewHandler(service)

		r.With(authMiddleware.RequireAuth).Get("/api/v1/profile", handler.GetProfile)
		r.With(authMiddleware.RequireAuth).Put("/api/v1/profile", handler.UpdateProfile)
	}
}

type ProfileResponse struct {
	User    model.User `json:"user"`
	Message string     `json:"message"`
}

func (handler *ProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetUserFromContext(r.Context())
	if !ok {
		handler.logger.Error("Failed to get user from context")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	ctx := r.Context()
	user, err := handler.service.GetProfile(ctx, claims.UserID)
	if err != nil {
		switch {
		case err.Error() == "user not found":
			http.Error(w, "User not found", http.StatusNotFound)
		case err.Error() == "user ID is required":
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			handler.logger.Error("Failed to get user profile", zap.Error(err))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	response := &ProfileResponse{
		User:    *user,
		Message: "Profile retrieved successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		handler.logger.Error("Failed to encode profile response", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	handler.logger.Info("Profile retrieved successfully", zap.String("user_id", claims.UserID))
}

type UpdateProfilePayload struct {
	Name     *string `json:"name"`
	Phone    *string `json:"phone"`
	Username *string `json:"username"`
}

func (handler *ProfileHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetUserFromContext(r.Context())
	if !ok {
		handler.logger.Error("Failed to get user from context")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var req UpdateProfilePayload
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		handler.logger.Error("Failed to decode request", zap.Error(err))
		return
	}
	defer r.Body.Close()

	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Phone != nil {
		updates["phone"] = *req.Phone
	}
	if req.Username != nil {
		updates["username"] = *req.Username
	}

	ctx := r.Context()
	user, err := handler.service.UpdateProfile(ctx, claims.UserID, updates)
	if err != nil {
		switch {
		case err.Error() == "user not found":
			http.Error(w, "User not found", http.StatusNotFound)
		case err.Error() == "username already exists":
			http.Error(w, "Username already exists", http.StatusConflict)
		case err.Error() == "user ID is required" || err.Error() == "no updates provided" || err.Error() == "no valid updates provided":
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			handler.logger.Error("Failed to update user profile", zap.Error(err))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	response := &ProfileResponse{
		User:    *user,
		Message: "Profile updated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		handler.logger.Error("Failed to encode profile response", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	handler.logger.Info("Profile updated successfully", zap.String("user_id", claims.UserID))
}
