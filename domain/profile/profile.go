package profile

import (
	"pokedex_backend_go/domain/profile/handler"
	"pokedex_backend_go/domain/profile/repository"
	"pokedex_backend_go/domain/profile/service"
	"pokedex_backend_go/pkg/auth"
	"pokedex_backend_go/pkg/server"

	"go.uber.org/fx"
)

func ProfileProvider() fx.Option {
	return fx.Options(
		fx.Provide(
			repository.NewRepository,
			service.NewService,
			auth.NewJWTService,
			auth.NewAuthMiddleware,
			server.AsHandler(handler.Handler),
			handler.NewHandler,
		),
	)
}
