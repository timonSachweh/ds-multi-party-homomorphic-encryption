package config

import (
	"context"
	"time"

	"github.com/sethvargo/go-envconfig"
)

// Configuration holds the configuration for the application.
type Configuration struct {
	HTTPServer
	PrivacyMLConfiguration
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
	PythonScriptPath      string `env:"PYTHON_SCRIPT_PATH,default=/python-ml/main.py"`
	MLModelPath           string `env:"ML_MODEL_PATH,default=/python-ml/model.pt"`
}

// Load processes the environment variables and returns the configuration.
func Load(ctx context.Context) (Configuration, error) {
	var cfg Configuration
	if err := envconfig.Process(ctx, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
