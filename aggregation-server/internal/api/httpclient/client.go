package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/aggregationserver/internal/config"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/aggregationserver/internal/entities"
)

type DataSpaceClientService interface {
	SendAggregatedResultsBack(body entities.MLModelWeights) error
}

type DataSpaceClientServiceImpl struct {
	ClientUrls []string
}

func NewDataSpaceClientService(clients config.ClientsConfiguration) DataSpaceClientService {
	return &DataSpaceClientServiceImpl{
		ClientUrls: clients.Urls,
	}
}

func (d *DataSpaceClientServiceImpl) SendAggregatedResultsBack(body entities.MLModelWeights) error {
	jsonData, err := json.Marshal(body)
	if err != nil {
		return err
	}
	_, err = http.Post(fmt.Sprintf("%s/v1/update", d.ClientUrls[0]), "application/json", bytes.NewReader(jsonData))
	return err
}
