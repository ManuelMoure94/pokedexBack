package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"pokedex_backend_go/pkg/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Service interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

func AsHandler(f interface{}) interface{} {
	return fx.Annotate(f, fx.ResultTags(`group:"handlers"`))
}

type params struct {
	fx.In
	Handlers []func(chi.Router) `group:"handlers"`
}

func New(params params) Service {
	srvLogger := logger.NewLogger("server")
	router := chi.NewRouter()

	srv := &service{
		logger: srvLogger,
		server: &http.Server{
			Addr:    fmt.Sprintf("%s:%s", "0.0.0.0", "3000"),
			Handler: router,
		},
	}

	serverTimeout, err := time.ParseDuration("60s")
	if err != nil {
		srv.logger.Error("Failed to parse server timeout", zap.Error(err))
	}

	router.Use(srv.RecoverMiddleware)
	router.Use(middleware.RealIP)
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.StripSlashes)
	router.Use(srv.CorsMiddleware)
	router.Use(middleware.Timeout(serverTimeout))

	router.Get("/ping", pingHandler)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	srv.logger.Info(fmt.Sprintf("Registering %d handlers", len(params.Handlers)))
	for _, handler := range params.Handlers {
		router.Group(handler)
	}

	return srv
}

type service struct {
	logger *zap.Logger
	server *http.Server
	pprof  *http.Server
}

func (s *service) Start(_ context.Context) error {
	pprofReady := make(chan struct{})

	go func() {
		router := chi.NewRouter()

		router.Use(s.RecoverMiddleware)
		router.Mount("/debug/pprof", middleware.Profiler())

		pprofHost := "localhost"
		pprofPort := "6060"

		s.pprof = &http.Server{
			Addr:    fmt.Sprintf("%s:%s", pprofHost, pprofPort),
			Handler: router,
		}

		s.logger.Info(fmt.Sprintf("Starting pprof server on http://%s", s.pprof.Addr))
		close(pprofReady)
		if err := s.pprof.ListenAndServe(); err != nil {
			s.logger.Error("Failed to start pprof server", zap.Error(err))
		}
	}()

	go func() {
		<-pprofReady
		s.logger.Info(fmt.Sprintf("Starting server on http://%s", s.server.Addr))
		if err := s.server.ListenAndServe(); err != nil {
			s.logger.Error("Failed to start server", zap.Error(err))
		}
	}()

	return nil
}

func (s *service) Stop(ctx context.Context) error {
	s.logger.Info("Stopping server")

	if s.pprof != nil {
		if err := s.pprof.Shutdown(ctx); err != nil {
			s.logger.Error("Failed to stop pprof server", zap.Error(err))
		}
	}

	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Error("Failed to stop server", zap.Error(err))
	}

	return nil
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "pong"}`))
}
