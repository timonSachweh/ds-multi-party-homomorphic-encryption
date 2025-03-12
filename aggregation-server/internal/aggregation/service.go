package aggregation

import (
	"github.com/robfig/cron/v3"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/aggregationserver/internal/api/httpclient"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/aggregationserver/internal/entities"
)

type AggregationService interface {
	AddNewData(entities.MLModelWeights)
	UpdateClients()
}

type AggregationServiceImpl struct {
	clientModelWeights []entities.MLModelWeights
	newData            bool
	httpClient         httpclient.DataSpaceClientService
}

// NewAggregationService creates a new instance of AggregationServiceImpl.
func NewAggregationService(httpClient httpclient.DataSpaceClientService) AggregationService {
	aggregationService := &AggregationServiceImpl{
		clientModelWeights: make([]entities.MLModelWeights, 0),
		newData:            false,
		httpClient:         httpClient,
	}
	c := cron.New()
	c.AddFunc("@every 00h00m10s", func() { aggregationService.UpdateClients() })
	c.Start()

	return aggregationService
}

func (a *AggregationServiceImpl) AddNewData(data entities.MLModelWeights) {
	a.clientModelWeights = append(a.clientModelWeights, data)
	a.newData = true
}

func (a *AggregationServiceImpl) UpdateClients() {
	if !a.newData {
		return
	}
	a.httpClient.SendAggregatedResultsBack(a.clientModelWeights[len(a.clientModelWeights)-1])
	a.newData = false
}
