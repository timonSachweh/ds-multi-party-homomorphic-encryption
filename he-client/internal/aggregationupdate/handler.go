package aggregationupdate

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/heclient/internal/entities"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/heclient/internal/ml"
)

// AggregationUpdateHandler defines the interface for handling aggregation updates.
type AggregationUpdateHandler interface {
	Routes() *chi.Mux
}

type aggregationUpdateHandlerImpl struct {
	mlService ml.MLService
	router    *chi.Mux
}

// NewAggregationUpdateHandler creates a new AggregationUpdateHandler with the provided service.
func NewAggregationUpdateHandler(mlService ml.MLService) AggregationUpdateHandler {
	return &aggregationUpdateHandlerImpl{
		mlService: mlService,
		router:    chi.NewRouter(),
	}
}

func (h *aggregationUpdateHandlerImpl) Routes() *chi.Mux {
	h.router.Post("/", h.handleUpdateModel)
	return h.router
}

func (h *aggregationUpdateHandlerImpl) handleUpdateModel(w http.ResponseWriter, r *http.Request) {
	var requestData entities.MLModelWeights
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	log.Println("Received model update request")
	h.mlService.UpdateModelWeights(requestData)
	w.WriteHeader(http.StatusOK)
}
