package main

import (
	"context"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/aggregationserver/internal/services"
	"log"

	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/aggregationserver/internal/api/http"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/aggregationserver/internal/api/httpclient"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/aggregationserver/internal/config"
)

func main() {
	ctx := context.Background()
	cfg, err := config.Load(ctx)
	if err != nil {
		log.Fatal(err)
	}

	httpClient := httpclient.NewDataSpaceClientService(cfg.PrivacyConfiguration)
	aggregationService := services.NewAggregationService(httpClient, cfg.PrivacyConfiguration)
	aggregationRouteHandler := http.NewAggregationRouteHandler(aggregationService)

	server := http.NewServer(cfg.HTTPServer, aggregationRouteHandler)
	server.Start(ctx)
}
