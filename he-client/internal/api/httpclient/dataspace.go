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
	UploadData(body entities.DataSpaceModelWeights) error
}

type DataSpaceClientServiceImpl struct {
	DataSpaceUrl string
}

func NewDataSpaceClientService(privacyMLConfig config.PrivacyMLConfiguration) DataSpaceClientService {
	return &DataSpaceClientServiceImpl{
		DataSpaceUrl: privacyMLConfig.AggregationServiceUrl,
	}
}

func (d *DataSpaceClientServiceImpl) UploadData(body entities.DataSpaceModelWeights) error {
	jsonData, err := json.Marshal(body)
	if err != nil {
		return err
	}
	_, err = http.Post(fmt.Sprintf("%s/clients/upload", d.DataSpaceUrl), "application/json", bytes.NewReader(jsonData))
	return err
}
