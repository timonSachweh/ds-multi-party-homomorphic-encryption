package http

import (
	"fmt"
	"net/http"
	"net/http/pprof"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func (s *Server) routes() {
	s.setupMiddlewares()
	s.setupDebugRoutes()

	v1 := chi.NewRouter()
	v1.Mount("/model", s.aggregationUpdateHandler.Routes())
	v1.Mount("/enc", s.encryptionHandler.Routes())
	s.router.Mount("/v1", v1)
	s.printRoutePaths()
}

func (s *Server) setupMiddlewares() {
	s.router.Use(middleware.Heartbeat("/health"))
	s.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	//s.router.Use(middleware.Logger)
}

func (s *Server) printRoutePaths() {
	chi.Walk(s.router, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fmt.Printf("[%s]: '%s' has %d middlewares\n", method, route, len(middlewares))
		return nil
	})
}

func (s *Server) setupDebugRoutes() {
	debug := chi.NewRouter()
	debug.Use(middleware.NoCache)
	debug.Use(middleware.StripSlashes)
	debug.HandleFunc("/pprof/*", pprof.Index)
	debug.HandleFunc("/pprof/cmdline", pprof.Cmdline)
	debug.HandleFunc("/pprof/profile", pprof.Profile)
	debug.HandleFunc("/pprof/symbol", pprof.Symbol)
	debug.HandleFunc("/pprof/trace", pprof.Trace)
	s.router.Mount("/debug", debug)
}
