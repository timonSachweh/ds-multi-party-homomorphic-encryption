package main

import (
	"context"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/heclient/internal/services"
	"log"

	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/heclient/internal/api/http"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/heclient/internal/api/httpclient"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/heclient/internal/config"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/heclient/internal/scheduling"
)

func main() {
	ctx := context.Background()
	cfg, err := config.Load(ctx)
	if err != nil {
		log.Fatal(err)
	}

	aggregationClient := httpclient.NewDataSpaceClientService(cfg.PrivacyMLConfiguration)
	pythonClient := httpclient.NewPythonClientService(cfg.PythonConfiguration)
	heService := services.NewHEService()
	mlService := services.NewMLService(heService, aggregationClient, pythonClient, cfg.PrivacyMLConfiguration)
	aggregationUpdateHandler := http.NewAggregationUpdateHandler(mlService)

	c := scheduling.InitializeCrons(mlService)
	defer c.Stop()

	server := http.NewServer(cfg.HTTPServer, aggregationUpdateHandler)
	server.Start(ctx)
}
