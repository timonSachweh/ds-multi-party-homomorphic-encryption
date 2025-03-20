package aggregation

import (
	"fmt"

	"github.com/robfig/cron/v3"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/aggregationserver/internal/api/httpclient"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/aggregationserver/internal/config"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/aggregationserver/internal/entities"
)

type AggregationService interface {
	AddNewData(entities.MLModelWeights)
	UpdateClients()
}

type AggregationServiceImpl struct {
	modelManager       map[string]ModelWeightStateManager
	minRequiredClients int
	httpClient         httpclient.DataSpaceClientService
}

// NewAggregationService creates a new instance of AggregationServiceImpl.
func NewAggregationService(httpClient httpclient.DataSpaceClientService, config config.PrivacyConfiguration) AggregationService {
	aggregationService := &AggregationServiceImpl{
		modelManager:       make(map[string]ModelWeightStateManager, 0),
		minRequiredClients: config.MinClientsNeeded,
		httpClient:         httpClient,
	}
	c := cron.New()
	c.AddFunc("@every 00h00m10s", func() { aggregationService.UpdateClients() })
	c.Start()

	return aggregationService
}

func (a *AggregationServiceImpl) AddNewData(data entities.MLModelWeights) {
	if data.ModelName == "" {
		return
	}
	if _, ok := a.modelManager[data.ModelName]; !ok {
		a.modelManager[data.ModelName] = NewModelWeightStateManager(a.minRequiredClients)
	}
	a.modelManager[data.ModelName].AddModelWeightsFromClient(data)
}

func (a *AggregationServiceImpl) UpdateClients() {
	fmt.Println("Updating clients")
	for _, modelManager := range a.modelManager {
		modelWeights, clients, err := modelManager.GetAggregatedModelWeights()
		if err != nil {
			fmt.Println(err)
			continue
		}
		a.httpClient.SendAggregatedResultsBack(clients, modelWeights)
	}
}
