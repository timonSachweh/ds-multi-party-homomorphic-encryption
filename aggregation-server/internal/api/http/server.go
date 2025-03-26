package http

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/aggregationserver/internal/config"
)

// Server represents an HTTP server with configuration and aggregation handler.
type Server struct {
	cfg                config.HTTPServer
	router             *chi.Mux
	aggregationHandler AggregationHandler
}

// NewServer creates a new HTTP server with the specified configuration and aggregation handler.
func NewServer(cfg config.HTTPServer, aggregationHandler AggregationHandler) *Server {
	srv := &Server{
		cfg:                cfg,
		router:             chi.NewRouter(),
		aggregationHandler: aggregationHandler,
	}

	srv.routes()
	return srv
}

// Start initializes and starts the HTTP server with the specified configuration.
// It listens on the port defined in the server configuration and handles incoming requests
// using the provided router. The server supports graceful shutdown, ensuring that all
// ongoing requests are completed before shutting down. The method logs the server's
// startup and shutdown events.
//
// Parameters:
//
//	ctx - the context to control the server's lifecycle and handle shutdown signals.
func (s *Server) Start(ctx context.Context) {
	server := http.Server{
		Addr:         fmt.Sprintf(":%d", s.cfg.Port),
		Handler:      s.router,
		IdleTimeout:  s.cfg.IdleTimeout,
		ReadTimeout:  s.cfg.ReadTimeout,
		WriteTimeout: s.cfg.WriteTimeout,
	}

	shutdownComplete := handleShutdown(func() {
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("server.Shutdown failed %v\n", err)
		}
	})

	log.Printf("Starting Webserver on Port: %d\n", s.cfg.Port)

	if err := server.ListenAndServe(); errors.Is(err, http.ErrServerClosed) {
		<-shutdownComplete
	} else {
		log.Printf("http.ListenAndServe failed: %v\n", err)
	}

	log.Println("Shutdown gracefully")
}

func handleShutdown(onShutdownSignal func()) <-chan struct{} {
	shutdown := make(chan struct{})

	go func() {
		shutdownSignal := make(chan os.Signal, 1)
		signal.Notify(shutdownSignal, os.Interrupt, syscall.SIGTERM)

		<-shutdownSignal

		onShutdownSignal()
		close(shutdown)
	}()
	return shutdown
}
