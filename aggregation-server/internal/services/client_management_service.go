package services

import (
	"github.com/Pro7ech/lattigo/rlwe"
	"github.com/robfig/cron/v3"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/aggregationserver/internal/api/httpclient"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/aggregationserver/internal/config"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/aggregationserver/internal/entities"
	"golang.org/x/crypto/openpgp/errors"
	"log"
	"time"
)

type ClientManagementService interface {
	AddNewData(entities.ClientModel) error
	UpdateClients()
	AddClient(clientModel entities.ClientModel) error
	GetClientsForModel(name string) ([]string, error)
	StartEncryptionSetupPhaseFor(modelName string)
	RequestClientTraining(modelName string)
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
	c.AddFunc("@every 00h00m30s", func() { aggregationService.UpdateClients() })
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
	for key, modelManager := range c.modelManager {
		if !modelManager.ModelCanAggregate() {
			continue
		}
		log.Println("Updating clients for model: " + key)
		_, modelWeights, clientUrls, weightLength := modelManager.GetClients()
		aggregatedModel := c.aggregateWeights(key, modelWeights, weightLength)
		c.initiateKeySwitchGeneration(clientUrls, aggregatedModel)
		c.encryptionService.CalculatePublicKeySwitchShare(modelManager.GetClientUrls())

		ciphertextWeights := aggregatedModel.WeightsAsCiphertext()
		ciphertextWeights = c.encryptionService.PublicKeySwitch(ciphertextWeights)

		weights := c.encryptionService.Decrypt(ciphertextWeights, weightLength)
		log.Printf("Weights decrypted with length: %d\n", len(weights))
		c.updateClientModels(modelManager.GetClientUrls(), weights, key)
		modelManager.ResetClientWeights()
	}
}

func (c *ClientManagementServiceImpl) aggregateWeights(key string, weights [][]*rlwe.Ciphertext, weightLength int) entities.ClientModel {
	updatedModelWeights := entities.ClientModel{
		ModelName: key,
		Length:    weightLength,
	}
	updatedModelWeights.SetCiphertextWeights(c.encryptionService.Aggregate(weights))
	return updatedModelWeights
}

func (c *ClientManagementServiceImpl) initiateKeySwitchGeneration(clientUrls []string, clientModel entities.ClientModel) {
	for _, client := range clientUrls {
		err := c.httpClient.SendPartialPublicKeySwitchGenerate(client, &clientModel)
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Println("Key switch generation initiated for clients")
}

func (c *ClientManagementServiceImpl) RequestClientTraining(modelName string) {
	if _, ok := c.modelManager[modelName]; !ok {
		log.Println("No such model and client combination is available")
		return
	}
	urls := c.modelManager[modelName].GetClientUrls()
	for _, url := range urls {
		err := c.httpClient.RequestClientTraining(url)
		if err != nil {
			log.Fatal(err)
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

	time.Sleep(5 * time.Second)

	clientUrls := c.modelManager[modelName].GetClientUrls()
	c.encryptionService.CalculatePublicKey(clientUrls)
	c.encryptionService.PublishPublicKey(clientUrls)
	c.encryptionService.CalculateRelinearizationKeys(clientUrls)
}

func (c *ClientManagementServiceImpl) updateClientModels(urls []string, weights []float64, key string) {
	modelUpdateResponse := entities.ModelClientUpdate{
		ModelName: key,
	}
	modelUpdateResponse.AddWeights(weights)
	for i := range urls {
		err := c.httpClient.SendAggregatedResultsBack(urls[i], &modelUpdateResponse)
		if err != nil {
			log.Println(err)
		}
	}
	log.Println("Updated client models for model: " + key)
}
