package ml

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/heclient/internal/api/httpclient"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/heclient/internal/config"
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
	config     config.PrivacyMLConfiguration
}

func NewMLService(heService privacy.HEService, httpClient httpclient.DataSpaceClientService, config config.PrivacyMLConfiguration) MLService {
	return &MLServiceImpl{
		model:      NewModel(),
		heService:  heService,
		httpClient: httpClient,
		config:     config,
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
	cmd := exec.Command("python3", m.config.PythonScriptPath, fmt.Sprintf("--model-path=%s", m.config.MLModelPath))
	out, err := cmd.Output()

	if err != nil {
		log.Println(err)
		return
	}

	log.Println(string(out))
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
