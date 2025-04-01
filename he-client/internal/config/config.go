package config

import (
	"context"
	"fmt"
	"time"

	"github.com/sethvargo/go-envconfig"
)

// Configuration holds the configuration for the application.
type Configuration struct {
	HTTPServer
	PrivacyMLConfiguration
	PythonConfiguration
}

// HTTPServer holds the configuration for the HTTP server.
type HTTPServer struct {
	IdleTimeout  time.Duration `env:"HTTP_IDLE_TIMEOUT,default=60s"`
	Port         int           `env:"HTTP_PORT,default=8080"`
	ReadTimeout  time.Duration `env:"HTTP_READ_TIMEOUT,default=1s"`
	WriteTimeout time.Duration `env:"HTTP_WRITE_TIMEOUT,default=2s"`
}

type PrivacyMLConfiguration struct {
	AggregationServiceUrl string `env:"AGGREGATION_SERVICE_URL,default=http://aggregation-server:8080"`
	ExternalUrl           string `env:"EXTERNAL_URL,default=http://localhost:8081"`
	ModelName             string `env:"MODEL_NAME,default=default_model"`
}

type PythonConfiguration struct {
	PythonHost    string `env:"PYTHON_HOST,default=http://localhost"`
	PythonPort    int    `env:"PYTHON_PORT,default=5000"`
	PythonApiPath string `env:"PYTHON_API_PATH,default=/api"`
}

func (c PythonConfiguration) BaseUrl() string {
	return fmt.Sprintf("%s:%d%s", c.PythonHost, c.PythonPort, c.PythonApiPath)
}

// Load processes the environment variables and returns the configuration.
func Load(ctx context.Context) (Configuration, error) {
	var cfg Configuration
	if err := envconfig.Process(ctx, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
