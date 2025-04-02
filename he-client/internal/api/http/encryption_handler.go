package http

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/heclient/internal/entities"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/heclient/internal/services"
	"log"
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
	h.router.Post("/gen/shared-public-key", h.handleSharedPublicKey)
	h.router.Put("/gen/shared-public-key", h.handleSharedPublicKey)
	h.router.Post("/gen/relinearization-key", h.handleSharedRelinKey)
	h.router.Put("/gen/relinearization-key", h.handleSharedRelinKey)
	h.router.Post("/public-key", h.handleReceivePublicKey)
	h.router.Put("/public-key", h.handleReceivePublicKey)
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

func (h *encryptionHandlerImpl) handleSharedRelinKey(w http.ResponseWriter, r *http.Request) {
	var share entities.RelinearizationKeyShare
	if err := json.NewDecoder(r.Body).Decode(&share); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	share = h.heService.PartialRelinKeyAggregation(share)

	render.JSON(w, r, share)
	w.WriteHeader(http.StatusOK)
}

func (h *encryptionHandlerImpl) handleReceivePublicKey(w http.ResponseWriter, r *http.Request) {
	var publicKey entities.PublicKeyExchange
	if err := json.NewDecoder(r.Body).Decode(&publicKey); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	err := h.heService.SetPublicKey(publicKey)
	if err != nil {
		log.Print(err)
	}
}
