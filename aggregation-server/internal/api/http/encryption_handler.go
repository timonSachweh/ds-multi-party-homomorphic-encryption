package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/aggregationserver/internal/services"
	"net/http"
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
	render.JSON(w, r, e.encryptionService.GetInformation())
}
