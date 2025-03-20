package config

import (
	"context"
	"time"

	"github.com/sethvargo/go-envconfig"
)

// Configuration holds the configuration for the HTTP server and the database.
// It embeds the HttpServer and Database structs.
type Configuration struct {
	HTTPServer
	PrivacyConfiguration
}

// HTTPServer holds the configuration for the HTTP server.
// The configuration values are loaded from environment variables.
type HTTPServer struct {
	IdleTimeout  time.Duration `env:"HTTP_IDLE_TIMEOUT,default=60s"` // IdleTimeout is the maximum amount of time to wait for the next request.
	Port         int           `env:"HTTP_PORT,default=8080"`        // Port is the port on which the server listens.
	ReadTimeout  time.Duration `env:"HTTP_READ_TIMEOUT,default=1s"`  // ReadTimeout is the maximum duration for reading the entire request.
	WriteTimeout time.Duration `env:"HTTP_WRITE_TIMEOUT,default=2s"` // WriteTimeout is the maximum duration before timing out writes of the response.
}

type PrivacyConfiguration struct {
	Urls             []string `env:"CLIENT_SERVICE_URLS,default=http://localhost:8081"`
	MinClientsNeeded int      `env:"MIN_CLIENTS_NEEDED,default=2"`
}

// Load loads the configuration from environment variables.
func Load(ctx context.Context) (Configuration, error) {
	var cfg Configuration
	if err := envconfig.Process(ctx, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
