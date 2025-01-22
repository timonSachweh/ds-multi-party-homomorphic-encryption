package ml

import (
	"log"

	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/heclient/internal/api/httpclient"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/heclient/internal/entities"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/heclient/internal/privacy"
)

type MLService interface {
	Train()
	Predict()
	RetrainAndSendUpdatedModelWeights()
	UpdateModelWeights(entities.MLModelWeights)
}

type MLServiceImpl struct {
	model      Model
	heService  privacy.HEService
	httpClient httpclient.DataSpaceClientService
}

func NewMLService(heService privacy.HEService, httpClient httpclient.DataSpaceClientService) MLService {
	return &MLServiceImpl{
		model:      NewModel(),
		heService:  heService,
		httpClient: httpClient,
	}
}

func (m *MLServiceImpl) RetrainAndSendUpdatedModelWeights() {
	m.Train()
	encrypt, err := m.heService.Encrypt(m.model.AsFloatVector())
	if err != nil {
		return
	}

	binary, err := encrypt.MarshalBinary()
	if err != nil {
		return
	}

	modelData := entities.MLModelWeights{
		ModelName: "model1",
		Weights:   binary,
		Length:    len(m.model.AsFloatVector()),
	}

	err = m.httpClient.UploadData(modelData)
	if err != nil {
		log.Fatal(err)
	}

}

func (m *MLServiceImpl) Train() {
	m.model.Train()
	log.Println(m.model.AsFloatVector())
}

func (m *MLServiceImpl) Predict() {
	m.model.Predict()
}

func (m *MLServiceImpl) UpdateModelWeights(weights entities.MLModelWeights) {
	decrypt, err := m.heService.Decrypt(weights.WeightsAsCiphertext(), weights.Length)
	if err != nil {
		return
	}

	m.model.UpdateWeights(decrypt)
}
