package register

import (
	"pokedex_backend_go/domain/register/handler"
	"pokedex_backend_go/domain/register/repository"
	"pokedex_backend_go/domain/register/service"
	"pokedex_backend_go/pkg/server"

	"go.uber.org/fx"
)

func RegisterProvider() fx.Option {
	return fx.Options(
		fx.Provide(
			repository.NewRepository,
			service.NewService,
			server.AsHandler(handler.Handler),
			handler.NewHandler,
		),
	)
}
