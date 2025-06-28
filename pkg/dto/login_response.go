package dto

import "pokedex_backend_go/pkg/model"

type LoginResponse struct {
	User  model.User `json:"user"`
	Token string     `json:"token"`
}
