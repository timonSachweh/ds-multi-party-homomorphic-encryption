package ml

import (
	"log"

	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/heclient/internal/api/httpclient"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/heclient/internal/config"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/heclient/internal/entities"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/heclient/internal/privacy"
)

type MLService interface {
	Train()
	Predict()
	RetrainAndSendUpdatedModelWeights()
	UpdateModelWeights(entities.DataSpaceModelWeights)
}

type MLServiceImpl struct {
	heService    privacy.HEService
	dspClient    httpclient.DataSpaceClientService
	pythonClient httpclient.PythonClientService
	config       config.PrivacyMLConfiguration
}

func NewMLService(heService privacy.HEService, dspClient httpclient.DataSpaceClientService, pythonClient httpclient.PythonClientService, config config.PrivacyMLConfiguration) MLService {
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

	modelData := entities.DataSpaceModelWeights{
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
	m.pythonClient.StartTraining()
}

func (m *MLServiceImpl) UpdateModelWeights(weights entities.DataSpaceModelWeights) {
	decrypt, err := m.heService.Decrypt(weights.WeightsAsCiphertext(), weights.Length)
	if err != nil {
		return
	}

	updatedModel := entities.MLModelWeights{
		ModelName: weights.ModelName,
		Weights:   decrypt,
		Length:    weights.Length,
	}

	m.pythonClient.UpdateModelWeights(updatedModel)
}
