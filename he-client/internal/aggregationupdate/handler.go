package aggregationupdate

import (
	"net/http"

	"github.com/go-chi/chi/v5"
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
	h.router.Post("/update", h.handleUpdateModel)
	return h.router
}

func (h *aggregationUpdateHandlerImpl) handleUpdateModel(w http.ResponseWriter, r *http.Request) {
}
