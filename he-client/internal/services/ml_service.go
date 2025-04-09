package services

import (
	"log"

	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/heclient/internal/api/httpclient"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/heclient/internal/config"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/heclient/internal/entities"
)

type MLService interface {
	Train()
	Predict()
	RetrainAndSendUpdatedModelWeights()
	UpdateModelWeights(update entities.ModelClientUpdate)
}

type MLServiceImpl struct {
	heService    HEService
	dspClient    httpclient.DataSpaceClientService
	pythonClient httpclient.PythonClientService
	config       config.PrivacyMLConfiguration
}

func NewMLService(heService HEService, dspClient httpclient.DataSpaceClientService, pythonClient httpclient.PythonClientService, config config.PrivacyMLConfiguration) MLService {
	return &MLServiceImpl{
		heService:    heService,
		dspClient:    dspClient,
		pythonClient: pythonClient,
		config:       config,
	}
}

func (m *MLServiceImpl) RetrainAndSendUpdatedModelWeights() {
	m.pythonClient.StartTraining()
	modelWeights, err := m.pythonClient.GetModelWeights()
	if err != nil {
		log.Fatal(err)
	}
	encrypt, err := m.heService.Encrypt(modelWeights.Weights)
	if err != nil {
		return
	}

	dataspaceModelWeightsPayload := make([][]byte, len(encrypt))
	for i, c := range encrypt {
		binary, err := c.MarshalBinary()
		if err != nil {
			return
		}
		dataspaceModelWeightsPayload[i] = binary
	}

	modelData := entities.ClientModel{
		ClientUrl: m.config.ExternalUrl,
		ModelName: modelWeights.ModelName,
		Weights:   dataspaceModelWeightsPayload,
		Length:    len(modelWeights.Weights),
	}

	err = m.dspClient.UploadData(modelData)
	if err != nil {
		log.Fatal(err)
	}

}

func (m *MLServiceImpl) Predict() {

}

func (m *MLServiceImpl) Train() {
	err := m.pythonClient.StartTraining()
	if err != nil {
		return
	}
}

func (m *MLServiceImpl) UpdateModelWeights(update entities.ModelClientUpdate) {
	updatedModel := entities.MLModelWeights{
		ModelName: update.ModelName,
		Weights:   update.Weights,
		Length:    update.Length,
	}

	err := m.pythonClient.UpdateModelWeights(updatedModel)
	if err != nil {
		return
	}
}
