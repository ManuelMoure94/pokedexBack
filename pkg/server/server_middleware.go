package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/cors"
	"go.uber.org/zap"
)

func (s *service) RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				if r == http.ErrAbortHandler {
					panic(r)
				}

				var err error
				switch x := r.(type) {
				case error:
					err = x
				default:
					err = fmt.Errorf("panic: %v", r)
				}

				s.logger.Error("Panic", zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (s *service) CorsMiddleware(next http.Handler) http.Handler {
	return cors.AllowAll().Handler(next)
}
