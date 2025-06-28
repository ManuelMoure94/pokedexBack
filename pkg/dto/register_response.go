package dto

import "pokedex_backend_go/pkg/model"

type RegisterResponse struct {
	User    model.User `json:"user"`
	Message string     `json:"message"`
	Token   string     `json:"token,omitempty"`
}
