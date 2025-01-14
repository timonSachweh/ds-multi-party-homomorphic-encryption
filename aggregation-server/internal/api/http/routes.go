package http

import (
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func (s *Server) routes() {
	s.setupMiddlewares()

	v1 := chi.NewRouter()
	v1.Mount("/agg", s.aggregationHandler.Routes())
	s.router.Mount("/v1", v1)
}

func (s *Server) setupMiddlewares() {
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Heartbeat("/health"))
	s.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
}
