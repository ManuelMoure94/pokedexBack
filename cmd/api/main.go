package main

import (
	"context"
	"fmt"
	"time"

	"pokedex_backend_go/domain/login"
	"pokedex_backend_go/domain/profile"
	"pokedex_backend_go/domain/register"
	"pokedex_backend_go/pkg/database"
	fxhelper "pokedex_backend_go/pkg/helper"
	"pokedex_backend_go/pkg/server"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

var logger *zap.Logger

func main() {
	logger = zap.L().WithOptions(zap.WithCaller(false)).Named("main")

	app := fx.New(
		// fx config
		fx.WithLogger(fxhelper.Logger),
		fx.StartTimeout(fxhelper.Timeout()),
		fx.StopTimeout(fxhelper.Timeout()),

		// Provide the database connection
		fx.Provide(database.Connection),
		fx.Provide(database.Gorm),
		fx.Invoke(database.Invoke),

		login.LoginProvider(),
		register.RegisterProvider(),
		profile.ProfileProvider(),

		fx.Provide(server.New),
		fx.Invoke(run),
	)

	defer func() {
		if r := recover(); r != nil {
			var err error
			switch x := r.(type) {
			case error:
				err = x
			default:
				err = fmt.Errorf("panic: %v", r)
			}

			logger.Error("Panic", zap.Error(err))
		}
	}()

	app.Run()
}

type runParams struct {
	fx.In
	Lifecycle fx.Lifecycle

	Server server.Service
}

func run(params runParams) {
	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return params.Server.Start(ctx)
		},
		OnStop: func(ctx context.Context) error {
			ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()

			return params.Server.Stop(ctx)
		},
	})
}
