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

	dataspaceClient := httpclient.NewDataSpaceClientService(cfg.PrivacyMLConfiguration)
	pythonClient := httpclient.NewPythonClientService(cfg.PythonConfiguration)
	heService := services.NewHEService()
	mlService := services.NewMLService(heService, dataspaceClient, pythonClient, cfg.PrivacyMLConfiguration)
	aggregationUpdateHandler := http.NewAggregationUpdateHandler(mlService)
	encryptionHandler := http.NewEncryptionHandler(heService)

	c := scheduling.InitializeCrons(mlService)
	defer c.Stop()

	go registerForDataSpaces(heService, dataspaceClient)

	server := http.NewServer(cfg.HTTPServer, aggregationUpdateHandler, encryptionHandler)
	server.Start(ctx)
}

func registerForDataSpaces(heService services.HEService, dataSpaceClient httpclient.DataSpaceClientService) {
	heConfiguration, err := dataSpaceClient.RegisterClient()
	if err != nil {
		log.Fatal(err)
		return
	}

	heService.SetParameters(heConfiguration)
}
