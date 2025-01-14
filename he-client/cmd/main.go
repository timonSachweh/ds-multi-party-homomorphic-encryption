package main

import (
	"context"
	"log"

	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/heclient/internal/aggregationupdate"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/heclient/internal/api/http"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/heclient/internal/config"
)

func main() {
	ctx := context.Background()
	cfg, err := config.Load(ctx)
	if err != nil {
		log.Fatal(err)
	}

	aggregationUpdateService := aggregationupdate.NewAggregationUpdateService()
	aggregationUpdateHandler := aggregationupdate.NewAggregationUpdateHandler(aggregationUpdateService)

	server := http.NewServer(cfg.HTTPServer, aggregationUpdateHandler)
	server.Start(ctx)
}
