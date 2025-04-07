package http

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/aggregationserver/internal/entities"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/aggregationserver/internal/services"
	"log"
	"net/http"
)

type ClientManagementHandler interface {
	Routes() *chi.Mux
}

type clientManagementRouteHandlerImpl struct {
	clientManagementService services.ClientManagementService
	encryptionService       services.EncryptionService
	router                  *chi.Mux
}

func NewClientManagementRouteHandler(clientManagementService services.ClientManagementService, encryptionService services.EncryptionService) ClientManagementHandler {
	return &clientManagementRouteHandlerImpl{
		clientManagementService: clientManagementService,
		encryptionService:       encryptionService,
		router:                  chi.NewRouter(),
	}
}

func (c *clientManagementRouteHandlerImpl) Routes() *chi.Mux {
	c.router.Post("/upload", c.handlePostClientData)
	c.router.Post("/", c.handleRegisterClient)
	c.router.Put("/", c.handleRegisterClient)
	c.router.Post("/train", c.handleTrainClients)
	c.router.Put("/train", c.handleTrainClients)
	c.router.Get("/train", c.handleTrainClients)
	return c.router
}

func (c *clientManagementRouteHandlerImpl) handlePostClientData(w http.ResponseWriter, r *http.Request) {
	var requestData entities.ClientModel
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	err := c.clientManagementService.AddNewData(requestData)
	if err != nil {
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (c *clientManagementRouteHandlerImpl) handleRegisterClient(w http.ResponseWriter, r *http.Request) {
	var requestData entities.ClientModel
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	err := c.clientManagementService.AddClient(requestData)
	if err != nil {
		log.Println(err)
	}

	go c.clientManagementService.StartEncryptionSetupPhaseFor(requestData.ModelName)

	render.JSON(w, r, c.encryptionService.GetInformation())
	w.WriteHeader(http.StatusOK)
}

func (c *clientManagementRouteHandlerImpl) handleTrainClients(w http.ResponseWriter, r *http.Request) {
	c.clientManagementService.RequestClientTraining("mnist")
}
