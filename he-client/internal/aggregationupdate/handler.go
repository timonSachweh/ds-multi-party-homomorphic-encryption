package aggregationupdate

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/heclient/internal/entities"
)

// AggregationUpdateHandler defines the interface for handling aggregation updates.
type AggregationUpdateHandler interface {
	Routes() *chi.Mux
}

type aggregationUpdateHandlerImpl struct {
	aggregationUpdateService AggregationUpdateService
	router                   *chi.Mux
}

// NewAggregationUpdateHandler creates a new AggregationUpdateHandler with the provided service.
func NewAggregationUpdateHandler(aggregationUpdateService AggregationUpdateService) AggregationUpdateHandler {
	return &aggregationUpdateHandlerImpl{
		aggregationUpdateService: aggregationUpdateService,
		router:                   chi.NewRouter(),
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
	log.Println(requestData.ModelName)
	w.WriteHeader(http.StatusOK)
}
