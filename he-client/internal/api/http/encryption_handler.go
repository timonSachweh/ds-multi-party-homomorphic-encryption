package http

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/heclient/internal/entities"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/heclient/internal/services"
	"net/http"
)

type EncryptionHandler interface {
	Routes() *chi.Mux
}

type encryptionHandlerImpl struct {
	router    *chi.Mux
	heService services.HEService
}

func NewEncryptionHandler(heService services.HEService) EncryptionHandler {
	return &encryptionHandlerImpl{
		router:    chi.NewRouter(),
		heService: heService,
	}
}

func (h *encryptionHandlerImpl) Routes() *chi.Mux {
	h.router.Post("/shared-public-key", h.handleSharedPublicKey)
	h.router.Put("/shared-public-key", h.handleSharedPublicKey)
	return h.router
}

func (h *encryptionHandlerImpl) handleSharedPublicKey(w http.ResponseWriter, r *http.Request) {
	var ckgShare entities.CkgShareExchange
	if err := json.NewDecoder(r.Body).Decode(&ckgShare); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	ckgShare = h.heService.PartialShareAggregation(ckgShare)

	render.JSON(w, r, ckgShare)
	w.WriteHeader(http.StatusOK)
}
