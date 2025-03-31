package http

import (
	"encoding/json"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/aggregationserver/internal/services"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/aggregationserver/internal/entities"
)

type ClientManagementHandler interface {
	Routes() *chi.Mux
}

type clientManagementRouteHandlerImpl struct {
	clientManagementService services.ClientManagementService
	router                  *chi.Mux
}

func NewAggregationRouteHandler(clientManagementService services.ClientManagementService) ClientManagementHandler {
	return &clientManagementRouteHandlerImpl{
		clientManagementService: clientManagementService,
		router:                  chi.NewRouter(),
	}
}

func (a *clientManagementRouteHandlerImpl) Routes() *chi.Mux {
	a.router.Post("/upload", a.handlePostClientData)
	return a.router
}

func (a *clientManagementRouteHandlerImpl) handlePostClientData(w http.ResponseWriter, r *http.Request) {
	var requestData entities.MLModelWeights
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	a.clientManagementService.AddNewData(requestData)
	w.WriteHeader(http.StatusOK)
}
