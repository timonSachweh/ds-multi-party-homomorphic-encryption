package aggregation

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type AggregationHandler interface {
	Routes() *chi.Mux
}

type aggregationRouteHandlerImpl struct {
	aggregationService AggregationService
	router             *chi.Mux
}

// NewAggregationRouteHandler creates a new AggregationHandler with the provided AggregationService.
func NewAggregationRouteHandler(aggregationService AggregationService) AggregationHandler {
	return &aggregationRouteHandlerImpl{
		aggregationService: aggregationService,
		router:             chi.NewRouter(),
	}
}

func (a *aggregationRouteHandlerImpl) Routes() *chi.Mux {
	a.router.Post("/", a.handlePostClientData)
	return a.router
}

func (a *aggregationRouteHandlerImpl) handlePostClientData(w http.ResponseWriter, r *http.Request) {
	var requestData MLModelWeights
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	fmt.Println(requestData)
}
