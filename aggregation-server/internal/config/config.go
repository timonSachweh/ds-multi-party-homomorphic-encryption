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
	Database
	ClientsConfiguration
}

// HTTPServer holds the configuration for the HTTP server.
// The configuration values are loaded from environment variables.
type HTTPServer struct {
	IdleTimeout  time.Duration `env:"HTTP_IDLE_TIMEOUT,default=60s"` // IdleTimeout is the maximum amount of time to wait for the next request.
	Port         int           `env:"HTTP_PORT,default=8080"`        // Port is the port on which the server listens.
	ReadTimeout  time.Duration `env:"HTTP_READ_TIMEOUT,default=1s"`  // ReadTimeout is the maximum duration for reading the entire request.
	WriteTimeout time.Duration `env:"HTTP_WRITE_TIMEOUT,default=2s"` // WriteTimeout is the maximum duration before timing out writes of the response.
}

// Database holds the configuration for the database connection.
// The configuration values are loaded from environment variables.
type Database struct {
	DatabaseURL                     string `env:"DB_URL,required"`                                        // DatabaseUrl is the URL of the database.
	Username                        string `env:"DB_USERNAME,required"`                                   // Username is the username for the database.
	Password                        string `env:"DB_PASSWORD,required"`                                   // Password is the password for the database.
	DatabaseName                    string `env:"DB_NAME,default=CV"`                                     // DatabaseName is the name of the database.
	AggregationServerCollectionName string `env:"DB_COLLECTION_NAME,default=AggregationServerCollection"` // BaseInformationCollectionName is the name of the base information collection.
}

type ClientsConfiguration struct {
	Urls []string `env:"CLIENT_SERVICE_URLS,default=http://localhost:8081"`
}

// Load loads the configuration from environment variables.
func Load(ctx context.Context) (Configuration, error) {
	var cfg Configuration
	if err := envconfig.Process(ctx, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
