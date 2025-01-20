package main

import (
	"context"
	"log"

	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/heclient/internal/aggregationupdate"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/heclient/internal/api/http"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/heclient/internal/api/httpclient"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/heclient/internal/config"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/heclient/internal/ml"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/heclient/internal/privacy"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/heclient/internal/scheduling"
)

func main() {
	ctx := context.Background()
	cfg, err := config.Load(ctx)
	if err != nil {
		log.Fatal(err)
	}

	httpClient := httpclient.NewDataSpaceClientService(cfg.PrivacyMLConfiguration)
	heService := privacy.NewHEService()
	mlService := ml.NewMLService(heService, httpClient)
	aggregationUpdateService := aggregationupdate.NewAggregationUpdateService()
	aggregationUpdateHandler := aggregationupdate.NewAggregationUpdateHandler(aggregationUpdateService)

	c := scheduling.InitializeCrons(mlService)
	defer c.Stop()

	server := http.NewServer(cfg.HTTPServer, aggregationUpdateHandler)
	server.Start(ctx)
}
