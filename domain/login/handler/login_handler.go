package handler

import (
	"encoding/json"
	"net/http"

	service "pokedex_backend_go/domain/login/service"
	"pokedex_backend_go/pkg/dto"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type LoginHandler struct {
	service *service.Service
	logger  *zap.Logger
}

func NewHandler(service *service.Service) *LoginHandler {
	return &LoginHandler{
		service: service,
		logger:  zap.L().Named("login_handler"),
	}
}

func Handler(service *service.Service) func(chi.Router) {
	return func(r chi.Router) {
		logger := zap.L().Named("login_handler_registration")
		logger.Info("Registering login handler at /api/v1/login")
		r.Post("/api/v1/login", NewHandler(service).LoginRequest)
	}
}

type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (handler *LoginHandler) LoginRequest(w http.ResponseWriter, r *http.Request) {
	var req LoginPayload
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		handler.logger.Error("Failed to decode request", zap.Error(err))
		return
	}
	defer r.Body.Close()

	if req.Email == "" || req.Password == "" {
		http.Error(w, "email and password are required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	user, token, err := handler.service.LoginWithToken(ctx, req.Email, req.Password)
	if err != nil {
		// Manejar diferentes tipos de errores
		switch {
		case err.Error() == "invalid username or password":
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		case err.Error() == "email is required" || err.Error() == "password is required":
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			handler.logger.Error("Failed to login user", zap.Error(err))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	response := &dto.LoginResponse{
		User:  *user,
		Token: token,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		handler.logger.Error("Failed to encode login response", zap.Error(err))
	}
	handler.logger.Info("Login successful", zap.String("email", req.Email))
}
