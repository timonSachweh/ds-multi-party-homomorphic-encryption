package http

import (
	"encoding/gob"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/aggregationserver/internal/services"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/aggregationserver/internal/utils"
)

type EncryptionHandler interface {
	Routes() *chi.Mux
}

type encryptionHandlerImpl struct {
	router            *chi.Mux
	encryptionService services.EncryptionService
}

func NewEncryptionHandler(encryptionService services.EncryptionService) EncryptionHandler {
	return &encryptionHandlerImpl{
		router:            chi.NewRouter(),
		encryptionService: encryptionService,
	}
}

func (e *encryptionHandlerImpl) Routes() *chi.Mux {
	e.router.Get("/information", e.handleGetInformation)
	return e.router
}

func (e *encryptionHandlerImpl) handleGetInformation(w http.ResponseWriter, r *http.Request) {
	utils.PrintMemoryStats("EncryptionHandler - handleGetInformation")
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(e.encryptionService.GetInformation())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	//render.JSON(w, r, e.encryptionService.GetInformation())
}
