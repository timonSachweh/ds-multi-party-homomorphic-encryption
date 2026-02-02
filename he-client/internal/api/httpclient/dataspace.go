package httpclient

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/heclient/internal/config"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/heclient/internal/entities"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/heclient/internal/utils"
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
	defer utils.PrintTime(time.Now(), "DataSpaceClientService - UploadData")
	jsonData, err := json.Marshal(body)
	if err != nil {
		return err
	}
	_, err = http.Post(fmt.Sprintf("%s/clients/upload", d.config.AggregationServiceUrl), "application/json", bytes.NewReader(jsonData))
	return err
}

func (d *DataSpaceClientServiceImpl) RegisterClient() (entities.PrivacyParams, error) {
	defer utils.PrintTime(time.Now(), "DataSpaceClientService - RegisterClient")
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
	decoder := gob.NewDecoder(resp.Body)
	var privacyParams entities.PrivacyParams
	err = decoder.Decode(&privacyParams)
	return privacyParams, err
}
