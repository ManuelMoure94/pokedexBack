package handler

import (
	"encoding/json"
	"net/http"

	"pokedex_backend_go/domain/register/service"
	"pokedex_backend_go/pkg/dto"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type RegisterHandler struct {
	service *service.Service
	logger  *zap.Logger
}

func NewHandler(service *service.Service) *RegisterHandler {
	return &RegisterHandler{
		service: service,
		logger:  zap.L().Named("register_handler"),
	}
}

func Handler(service *service.Service) func(chi.Router) {
	return func(r chi.Router) {
		logger := zap.L().Named("register_handler_registration")
		logger.Info("Registering register handler at /api/v1/register")
		r.Post("/api/v1/register", NewHandler(service).RegisterRequest)
	}
}

type RegisterPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (handler *RegisterHandler) RegisterRequest(w http.ResponseWriter, r *http.Request) {
	var req RegisterPayload
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		handler.logger.Error("Failed to decode request", zap.Error(err))
		return
	}
	defer r.Body.Close()

	if req.Email == "" || req.Password == "" {
		http.Error(w, "email and password are required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	user, token, err := handler.service.RegisterWithToken(ctx, req.Email, req.Password)
	if err != nil {
		switch {
		case err.Error() == "email already exists":
			http.Error(w, "Email already exists", http.StatusConflict)
		case err.Error() == "invalid email format":
			http.Error(w, "Invalid email format", http.StatusBadRequest)
		case err.Error() == "password must be at least 6 characters long":
			http.Error(w, "Password must be at least 6 characters long", http.StatusBadRequest)
		case err.Error() == "email is required" || err.Error() == "password is required":
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			handler.logger.Error("Failed to register user", zap.Error(err))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	response := &dto.RegisterResponse{
		User:    *user,
		Message: "User registered successfully",
		Token:   token,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		handler.logger.Error("Failed to encode response", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	handler.logger.Info("User registered successfully", zap.String("email", req.Email))
}
