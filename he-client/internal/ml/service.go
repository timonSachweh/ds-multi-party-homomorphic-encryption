package ml

import (
	"log"

	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/heclient/internal/api/httpclient"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/heclient/internal/privacy"
)

type MLService interface {
	Train()
	Predict()
	RetrainAndSendUpdatedModelWeights()
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

	modelData := httpclient.MLModelWeights{
		ModelName: "model1",
		Weights:   *encrypt,
		Length:    len(m.model.AsFloatVector()),
	}

	err = m.httpClient.UploadData(modelData)
	if err != nil {
		log.Fatal(err)
	}

}

func (m *MLServiceImpl) Train() {
	m.model.Train()
}

func (m *MLServiceImpl) Predict() {
	m.model.Predict()
}
