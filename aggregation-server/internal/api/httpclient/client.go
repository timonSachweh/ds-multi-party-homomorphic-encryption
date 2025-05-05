package httpclient

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Pro7ech/lattigo/mhe"
	"io"
	"net/http"

	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/aggregationserver/internal/entities"
)

type DataSpaceClientService interface {
	SendAggregatedResultsBack(url string, body *entities.ModelClientUpdate) error
	SendPartialPublicCkgShare(url string, body *entities.CkgShareExchange) error
	SendPublicKeyToClient(url string, pk *entities.PublicKeyExchange) error
	SendPartialRelinearizationKey(url string, e *entities.RelinearizationKeyShare) error
	SendPartialPublicKeySwitchGenerate(url string, e *entities.ClientModel) error
	SendPartialPublicKeySwitchAggregation(url string, e *mhe.KeySwitchingShare) error
	RequestClientTraining(url string) error
}

type DataSpaceClientServiceImpl struct {
}

func NewDataSpaceClientService() DataSpaceClientService {
	return &DataSpaceClientServiceImpl{}
}

func (d *DataSpaceClientServiceImpl) RequestClientTraining(url string) error {
	if url == "" {
		return errors.New("url is empty")
	}
	_, err := http.Get(fmt.Sprintf("%s/v1/model/train", url))
	return err
}

func (d *DataSpaceClientServiceImpl) SendAggregatedResultsBack(url string, body *entities.ModelClientUpdate) error {
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
	var data bytes.Buffer
	encoder := gob.NewEncoder(&data)
	err := encoder.Encode(body)
	if err != nil {
		return err
	}
	if url == "" {
		return errors.New("url is empty")
	}
	resp, err := http.Post(fmt.Sprintf("%s/v1/enc/gen/shared-public-key", url), "application/json", &data)
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(body)
	return err
}

func (d *DataSpaceClientServiceImpl) SendPartialRelinearizationKey(url string, e *entities.RelinearizationKeyShare) error {
	var data bytes.Buffer
	encoder := gob.NewEncoder(&data)
	err := encoder.Encode(e)
	if err != nil {
		return err
	}
	if url == "" {
		return errors.New("url is empty")
	}
	resp, err := http.Post(fmt.Sprintf("%s/v1/enc/gen/relinearization-key", url), "application/json", &data)
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(e)
	if err != nil {
		return err
	}
	return resp.Body.Close()
}

func (d *DataSpaceClientServiceImpl) SendPartialPublicKeySwitchGenerate(url string, e *entities.ClientModel) error {
	var data bytes.Buffer
	encoder := gob.NewEncoder(&data)
	err := encoder.Encode(e)
	if err != nil {
		return err
	}
	if url == "" {
		return errors.New("url is empty")
	}
	_, err = http.Post(fmt.Sprintf("%s/v1/enc/gen/public-key-switch", url), "application/json", &data)
	return err
}

func (d *DataSpaceClientServiceImpl) SendPartialPublicKeySwitchAggregation(url string, e *mhe.KeySwitchingShare) error {
	binaryData, err := e.MarshalBinary()
	if err != nil {
		return err
	}
	if url == "" {
		return errors.New("url is empty")
	}
	resp, err := http.Post(fmt.Sprintf("%s/v1/enc/gen/public-key-switch-aggregate", url), "application/json", bytes.NewReader(binaryData))
	if err != nil {
		return err
	}
	byteValue, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = e.UnmarshalBinary(byteValue)
	return resp.Body.Close()
}

func (d *DataSpaceClientServiceImpl) SendPublicKeyToClient(url string, pk *entities.PublicKeyExchange) error {
	var data bytes.Buffer
	encoder := gob.NewEncoder(&data)
	err := encoder.Encode(pk)
	if err != nil {
		return err
	}
	if url == "" {
		return errors.New("url is empty")
	}
	resp, err := http.Post(fmt.Sprintf("%s/v1/enc/public-key", url), "application/json", &data)
	if err != nil {
		return err
	}
	return resp.Body.Close()
}
