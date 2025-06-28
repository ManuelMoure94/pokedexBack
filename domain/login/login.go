package login

import (
	"pokedex_backend_go/domain/login/handler"
	"pokedex_backend_go/domain/login/repository"
	"pokedex_backend_go/domain/login/service"
	"pokedex_backend_go/pkg/server"

	"go.uber.org/fx"
)

func LoginProvider() fx.Option {
	return fx.Options(
		fx.Provide(
			repository.NewRepository,
			service.NewService,
			server.AsHandler(handler.Handler),
			handler.NewHandler,
		),
	)
}
