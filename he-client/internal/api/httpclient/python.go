package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/heclient/internal/config"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/heclient/internal/entities"
)

type PythonClientService interface {
	Predict(body entities.PredictionRequest) (entities.PredictionResponse, error)
	GetModelWeights() (entities.MLModelWeights, error)
	UpdateModelWeights(entities.MLModelWeights) error
	StartTraining() error
}

type PythonClientServiceImpl struct {
	PythonServiceUrl string
}

func NewPythonClientService(pythonConfig config.PythonConfiguration) PythonClientService {
	return &PythonClientServiceImpl{
		PythonServiceUrl: pythonConfig.BaseUrl(),
	}
}

func (p *PythonClientServiceImpl) Predict(body entities.PredictionRequest) (entities.PredictionResponse, error) {
	jsonData, err := json.Marshal(body)
	if err != nil {
		return entities.PredictionResponse{}, err
	}
	resp, err := http.Post(fmt.Sprintf("%s/predict", p.PythonServiceUrl), "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return entities.PredictionResponse{}, err
	}
	defer resp.Body.Close()
	var prediction entities.PredictionResponse
	err = json.NewDecoder(resp.Body).Decode(&prediction)
	return prediction, err
}

func (p *PythonClientServiceImpl) GetModelWeights() (entities.MLModelWeights, error) {
	resp, err := http.Get(fmt.Sprintf("%s/model-params", p.PythonServiceUrl))
	if err != nil {
		return entities.MLModelWeights{}, err
	}
	defer resp.Body.Close()
	var modelWeights entities.MLModelWeights
	err = json.NewDecoder(resp.Body).Decode(&modelWeights)
	return modelWeights, err
}

func (p *PythonClientServiceImpl) UpdateModelWeights(modelWeights entities.MLModelWeights) error {
	jsonData, err := json.Marshal(modelWeights)
	if err != nil {
		fmt.Println(err)
		return err
	}
	_, err = http.Post(fmt.Sprintf("%s/model-params", p.PythonServiceUrl), "application/json", bytes.NewReader(jsonData))
	return err
}

func (p *PythonClientServiceImpl) StartTraining() error {
	_, err := http.Get(fmt.Sprintf("%s/train", p.PythonServiceUrl))
	return err
}
