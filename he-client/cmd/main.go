package main

import (
	"context"
	"log"

	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/heclient/internal/aggregationupdate"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/heclient/internal/api/http"
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

	mlService := ml.NewMLService()
	heService := privacy.NewHEService()
	aggregationUpdateService := aggregationupdate.NewAggregationUpdateService()
	aggregationUpdateHandler := aggregationupdate.NewAggregationUpdateHandler(aggregationUpdateService)

	c := scheduling.InitializeCrons(mlService, heService)
	defer c.Stop()

	server := http.NewServer(cfg.HTTPServer, aggregationUpdateHandler)
	server.Start(ctx)
}
