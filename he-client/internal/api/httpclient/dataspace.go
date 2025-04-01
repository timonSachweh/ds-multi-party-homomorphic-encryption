package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/heclient/internal/config"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/heclient/internal/entities"
)

type DataSpaceClientService interface {
	UploadData(body entities.ClientModel) error
	RegisterClient() (entities.PrivacyParams, error)
}

type DataSpaceClientServiceImpl struct {
	config config.PrivacyMLConfiguration
}

func NewDataSpaceClientService(privacyMLConfig config.PrivacyMLConfiguration) DataSpaceClientService {
	return &DataSpaceClientServiceImpl{
		config: privacyMLConfig,
	}
}

func (d *DataSpaceClientServiceImpl) UploadData(body entities.ClientModel) error {
	jsonData, err := json.Marshal(body)
	if err != nil {
		return err
	}
	_, err = http.Post(fmt.Sprintf("%s/clients/upload", d.config.AggregationServiceUrl), "application/json", bytes.NewReader(jsonData))
	return err
}

func (d *DataSpaceClientServiceImpl) RegisterClient() (entities.PrivacyParams, error) {
	model := entities.ClientModel{
		ClientUrl: d.config.ExternalUrl,
		ModelName: d.config.ModelName,
	}

	jsonData, err := json.Marshal(model)
	if err != nil {
		return entities.PrivacyParams{}, err
	}
	resp, err := http.Post(fmt.Sprintf("%s/clients", d.config.AggregationServiceUrl), "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return entities.PrivacyParams{}, err
	}
	decoder := json.NewDecoder(resp.Body)
	var privacyParams entities.PrivacyParams
	err = decoder.Decode(&privacyParams)
	return privacyParams, err
}
