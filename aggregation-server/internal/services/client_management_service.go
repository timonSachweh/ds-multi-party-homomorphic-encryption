package services

import (
	"github.com/robfig/cron/v3"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/aggregationserver/internal/api/httpclient"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/aggregationserver/internal/config"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/aggregationserver/internal/entities"
)

type ClientManagementService interface {
	AddNewData(entities.MLModelWeights)
	UpdateClients()
}

type ClientManagementServiceImpl struct {
	modelManager       map[string]ModelWeightStateManager
	minRequiredClients int
	httpClient         httpclient.DataSpaceClientService
}

// NewAggregationService creates a new instance of AggregationServiceImpl.
func NewClientManagementService(httpClient httpclient.DataSpaceClientService, config config.PrivacyConfiguration) ClientManagementService {
	aggregationService := &ClientManagementServiceImpl{
		modelManager:       make(map[string]ModelWeightStateManager, 0),
		minRequiredClients: config.MinClientsNeeded,
		httpClient:         httpClient,
	}
	c := cron.New()
	c.AddFunc("@every 00h00m10s", func() { aggregationService.UpdateClients() })
	c.Start()

	return aggregationService
}

func (a *ClientManagementServiceImpl) AddNewData(data entities.MLModelWeights) {
	if data.ModelName == "" {
		return
	}
	if _, ok := a.modelManager[data.ModelName]; !ok {
		a.modelManager[data.ModelName] = NewModelWeightStateManager(a.minRequiredClients)
	}
	a.modelManager[data.ModelName].AddModelWeightsFromClient(data)
}

func (a *ClientManagementServiceImpl) UpdateClients() {
	for _, modelManager := range a.modelManager {
		modelWeights, clients, err := modelManager.GetAggregatedModelWeights()
		if err != nil {
			// will be the error for "not enough clients to aggregate"
			continue
		}
		a.httpClient.SendAggregatedResultsBack(clients, modelWeights)
	}
}
