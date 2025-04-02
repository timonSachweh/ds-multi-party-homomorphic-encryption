package services

import (
	"github.com/robfig/cron/v3"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/aggregationserver/internal/api/httpclient"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/aggregationserver/internal/config"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/aggregationserver/internal/entities"
	"golang.org/x/crypto/openpgp/errors"
	"log"
)

type ClientManagementService interface {
	AddNewData(entities.ClientModel) error
	UpdateClients()
	AddClient(clientModel entities.ClientModel) error
	GetClientsForModel(name string) ([]string, error)
	StartEncryptionSetupPhaseFor(modelName string)
}

type ClientManagementServiceImpl struct {
	modelManager       map[string]ModelWeightStateManager
	minRequiredClients int
	httpClient         httpclient.DataSpaceClientService
	encryptionService  EncryptionService
}

func NewClientManagementService(httpClient httpclient.DataSpaceClientService, encryptionService EncryptionService, config config.PrivacyConfiguration) ClientManagementService {
	aggregationService := &ClientManagementServiceImpl{
		modelManager:       make(map[string]ModelWeightStateManager, 0),
		minRequiredClients: config.MinClientsNeeded,
		httpClient:         httpClient,
		encryptionService:  encryptionService,
	}
	c := cron.New()
	c.AddFunc("@every 00h00m10s", func() { aggregationService.UpdateClients() })
	c.Start()

	return aggregationService
}

func (c *ClientManagementServiceImpl) AddNewData(data entities.ClientModel) error {
	if data.ModelName == "" {
		return errors.InvalidArgumentError("model name is required")
	}
	if _, ok := c.modelManager[data.ModelName]; !ok {
		return errors.InvalidArgumentError("No such model and client combination is available")
	}
	c.modelManager[data.ModelName].AddModelWeightsFromClient(data)
	return nil
}

func (c *ClientManagementServiceImpl) UpdateClients() {
	for _, modelManager := range c.modelManager {
		modelWeights, clients, err := modelManager.GetAggregatedModelWeights()
		if err != nil {
			// will be the error for "not enough clients to aggregate"
			continue
		}
		for i := range clients {
			err = c.httpClient.SendAggregatedResultsBack(clients[i], modelWeights)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func (c *ClientManagementServiceImpl) AddClient(clientModel entities.ClientModel) error {
	if clientModel.ModelName == "" || clientModel.ClientUrl == "" {
		return errors.InvalidArgumentError("Client model name or client url is empty")
	}
	if _, ok := c.modelManager[clientModel.ModelName]; !ok {
		c.modelManager[clientModel.ModelName] = NewModelWeightStateManager(c.minRequiredClients)
	}
	err := c.modelManager[clientModel.ModelName].AddClient(clientModel)
	return err
}

func (c *ClientManagementServiceImpl) GetClientsForModel(name string) ([]string, error) {
	if name == "" {
		return nil, errors.InvalidArgumentError("model name is required")
	}
	if _, ok := c.modelManager[name]; !ok {
		return nil, errors.InvalidArgumentError("No such model and client combination is available")
	}
	return c.modelManager[name].GetClientUrls(), nil
}

func (c *ClientManagementServiceImpl) StartEncryptionSetupPhaseFor(modelName string) {
	if c.modelManager[modelName] == nil {
		return
	}

	clientUrls := c.modelManager[modelName].GetClientUrls()
	c.encryptionService.CalculatePublicKey(clientUrls)
	c.encryptionService.PublishPublicKey(clientUrls)
	c.encryptionService.CalculateRelinearizationKeys(clientUrls)

}
