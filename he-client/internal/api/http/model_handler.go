package http

import (
	"encoding/json"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/heclient/internal/services"
	"io"
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
	mlService services.MLService
	router    *chi.Mux
}

// NewAggregationUpdateHandler creates a new AggregationUpdateHandler with the provided service.
func NewAggregationUpdateHandler(mlService services.MLService) AggregationUpdateHandler {
	return &aggregationUpdateHandlerImpl{
		mlService: mlService,
		router:    chi.NewRouter(),
	}
}

func (h *aggregationUpdateHandlerImpl) Routes() *chi.Mux {
	h.router.Post("/updated-model", h.handleUpdateModel)
	h.router.Get("/train", h.handleTrainModel)
	h.router.Post("/train", h.handleTrainModel)
	h.router.Put("/train", h.handleTrainModel)
	h.router.Post("/predict", h.handlePrediction)
	h.router.Put("/predict", h.handlePrediction)
	h.router.Post("/predict-image", h.handlePredictionImage)
	h.router.Put("/predict-image", h.handlePredictionImage)
	return h.router
}

func (h *aggregationUpdateHandlerImpl) handleUpdateModel(w http.ResponseWriter, r *http.Request) {
	var requestData entities.ClientModel
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Println(err)
		return
	}
	log.Println("Handler: model update request")
	h.mlService.UpdateModelWeights(requestData)
	w.WriteHeader(http.StatusOK)
}

func (h *aggregationUpdateHandlerImpl) handleTrainModel(w http.ResponseWriter, r *http.Request) {
	log.Println("Handler: model training request")
	go h.mlService.RetrainAndSendUpdatedModelWeights()
	w.WriteHeader(http.StatusOK)
}

func (h *aggregationUpdateHandlerImpl) handlePrediction(w http.ResponseWriter, r *http.Request) {
	log.Println("Handler: model prediction request")
	var requestData entities.PredictionRequest
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		log.Fatal(err)
	}
	h.mlService.Predict()
	w.WriteHeader(http.StatusOK)
}

func (h *aggregationUpdateHandlerImpl) handlePredictionImage(w http.ResponseWriter, r *http.Request) {
	log.Println("Handler: model image prediction request")
	if err := r.ParseMultipartForm(8 * 1024 * 1024); err != nil { // 8 MB
		log.Printf("Could not parse multipart form: %v\n", err)
		return
	}
	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		return
	}
	defer file.Close()

	fileSize := fileHeader.Size
	if fileSize > 8*1024*1024 {
		return
	}
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return
	}

	detectedFileType := http.DetectContentType(fileBytes)
	switch detectedFileType {
	case "image/jpeg", "image/jpg", "image/png":
		break
	default:
		return
	}
	// h.mlService.PredictImage()
	w.WriteHeader(http.StatusOK)
}
