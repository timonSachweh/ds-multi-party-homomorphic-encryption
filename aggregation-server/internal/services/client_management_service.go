package services

import (
	"log"
	"time"

	"github.com/Pro7ech/lattigo/rlwe"
	"github.com/robfig/cron/v3"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/aggregationserver/internal/api/httpclient"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/aggregationserver/internal/config"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/aggregationserver/internal/entities"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/aggregationserver/internal/utils"
	"golang.org/x/crypto/openpgp/errors"
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
	c.AddFunc("@every 00h00m5s", func() { aggregationService.UpdateClients() })
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
	utils.PrintMemoryStats("UpdateClients - Initiated")
	for key, modelManager := range c.modelManager {
		if !modelManager.ModelCanAggregate() {
			continue
		}
		log.Println("Updating clients for model: " + key)
		_, modelWeights, clientUrls, weightLength := modelManager.GetClients()
		utils.PrintMemoryStats("UpdateClients - BeforeAggregation")
		aggregatedModel := c.aggregateWeights(key, modelWeights, weightLength)
		utils.PrintMemoryStats("UpdateClients - AfterAggregation")
		c.initiateKeySwitchGeneration(clientUrls, aggregatedModel)
		utils.PrintMemoryStats("UpdateClients - AfterKeySwitchGeneration")
		c.encryptionService.CalculatePublicKeySwitchShare(modelManager.GetClientUrls())
		utils.PrintMemoryStats("UpdateClients - AfterPublicKeySwitchShareCalculation")

		ciphertextWeights := aggregatedModel.WeightsAsCiphertext()
		ciphertextWeights = c.encryptionService.PublicKeySwitch(ciphertextWeights)

		utils.PrintMemoryStats("UpdateClients - AfterPublicKeySwitch")

		weights := c.encryptionService.Decrypt(ciphertextWeights, weightLength)
		log.Printf("Weights decrypted with length: %d\n", len(weights))
		c.updateClientModels(modelManager.GetClientUrls(), weights, key)
		modelManager.ResetClientWeights()

		utils.PrintMemoryStats("UpdateClients - AfterUpdateClients")
	}
}

func (c *ClientManagementServiceImpl) aggregateWeights(key string, weights [][]*rlwe.Ciphertext, weightLength int) entities.ClientModel {
	defer utils.PrintTime(time.Now(), "aggregateWeights")
	updatedModelWeights := entities.ClientModel{
		ModelName: key,
		Length:    weightLength,
	}
	updatedModelWeights.SetCiphertextWeights(c.encryptionService.Aggregate(weights))
	return updatedModelWeights
}

func (c *ClientManagementServiceImpl) initiateKeySwitchGeneration(clientUrls []string, clientModel entities.ClientModel) {
	defer utils.PrintTime(time.Now(), "initiateKeySwitchGeneration")
	for _, client := range clientUrls {
		err := c.httpClient.SendPartialPublicKeySwitchGenerate(client, &clientModel)
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Println("Key switch generation initiated for clients")
}

func (c *ClientManagementServiceImpl) RequestClientTraining(modelName string) {
	defer utils.PrintTime(time.Now(), "requestClientTraining")
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
	defer utils.PrintTime(time.Now(), "addClient")
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
	defer utils.PrintTime(time.Now(), "getClientsForModel")
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
	defer utils.PrintTime(time.Now(), "startEncryptionSetupPhaseFor")

	time.Sleep(5 * time.Second)
	utils.PrintMemoryStats("ClientManagementService - EncryptionSetupBegin")
	clientUrls := c.modelManager[modelName].GetClientUrls()
	c.encryptionService.CalculatePublicKey(clientUrls)
	utils.PrintMemoryStats("ClientManagementService - EncryptionSetupPublicKey")
	c.encryptionService.PublishPublicKey(clientUrls)
	utils.PrintMemoryStats("ClientManagementService - EncryptionSetupPublicKeyPushed")
	c.encryptionService.CalculateRelinearizationKeys(clientUrls)
	utils.PrintMemoryStats("ClientManagementService - EncryptionSetupRelinearizationKeys")
}

func (c *ClientManagementServiceImpl) updateClientModels(urls []string, weights []float64, key string) {
	defer utils.PrintTime(time.Now(), "updateClientModels")
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
