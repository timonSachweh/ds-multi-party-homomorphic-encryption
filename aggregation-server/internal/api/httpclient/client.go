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
	SendAggregatedResultsBack(clients []string, body entities.MLModelWeights) error
}

type DataSpaceClientServiceImpl struct {
	ClientUrls []string
}

func NewDataSpaceClientService(clients config.PrivacyConfiguration) DataSpaceClientService {
	return &DataSpaceClientServiceImpl{
		ClientUrls: clients.Urls,
	}
}

func (d *DataSpaceClientServiceImpl) SendAggregatedResultsBack(clients []string, body entities.MLModelWeights) error {
	jsonData, err := json.Marshal(body)
	if err != nil {
		return err
	}
	if clients == nil {
		clients = d.ClientUrls
	}
	for _, url := range clients {
		fmt.Println("Sending data to: ", url)
		_, err = http.Post(fmt.Sprintf("%s/v1/updated-model", url), "application/json", bytes.NewReader(jsonData))
		if err != nil {
			return err
		}
	}
	return err
}
