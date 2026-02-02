package main

import (
	"context"
	"log"

	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/aggregationserver/internal/services"

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
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	httpClient := httpclient.NewDataSpaceClientService()
	encryptionService := services.NewEncryptionService(httpClient)
	clientManagementService := services.NewClientManagementService(httpClient, encryptionService, cfg.PrivacyConfiguration)
	clientManagementRouteHandler := http.NewClientManagementRouteHandler(clientManagementService, encryptionService)
	encryptionRouteHandler := http.NewEncryptionHandler(encryptionService)

	server := http.NewServer(cfg.HTTPServer, clientManagementRouteHandler, encryptionRouteHandler)
	server.Start(ctx)
}
