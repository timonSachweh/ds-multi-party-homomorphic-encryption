package httpclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/aggregationserver/internal/entities"
)

type DataSpaceClientService interface {
	SendAggregatedResultsBack(url string, body entities.ClientModel) error
	SendPartialPublicCkgShare(url string, body *entities.CkgShareExchange) error
}

type DataSpaceClientServiceImpl struct {
}

func NewDataSpaceClientService() DataSpaceClientService {
	return &DataSpaceClientServiceImpl{}
}

func (d *DataSpaceClientServiceImpl) SendAggregatedResultsBack(url string, body entities.ClientModel) error {
	jsonData, err := json.Marshal(body)
	if err != nil {
		return err
	}
	if url == "" {
		return errors.New("url is empty")
	}
	_, err = http.Post(fmt.Sprintf("%s/v1/model/updated-model", url), "application/json", bytes.NewReader(jsonData))
	return err
}

func (d *DataSpaceClientServiceImpl) SendPartialPublicCkgShare(url string, body *entities.CkgShareExchange) error {
	jsonData, err := json.Marshal(body)
	if err != nil {
		return err
	}
	if url == "" {
		return errors.New("url is empty")
	}
	resp, err := http.Post(fmt.Sprintf("%s/v1/enc/shared-public-key", url), "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(body)
	return err
}
