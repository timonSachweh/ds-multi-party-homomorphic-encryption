package aggregation

import (
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/aggregationserver/internal/api/httpclient"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/aggregationserver/internal/entities"
)

type AggregationService interface {
	AddNewData(entities.MLModelWeights)
	UpdateClients()
}

type AggregationServiceImpl struct {
	clientModelWeights []entities.MLModelWeights
	httpClient         httpclient.DataSpaceClientService
}

// NewAggregationService creates a new instance of AggregationServiceImpl.
func NewAggregationService(httpClient httpclient.DataSpaceClientService) AggregationService {
	return &AggregationServiceImpl{
		clientModelWeights: make([]entities.MLModelWeights, 0),
		httpClient:         httpClient,
	}
}

func (a *AggregationServiceImpl) AddNewData(data entities.MLModelWeights) {
	a.clientModelWeights = append(a.clientModelWeights, data)
	a.UpdateClients()

}

func (a *AggregationServiceImpl) UpdateClients() {
	a.httpClient.SendAggregatedResultsBack(a.clientModelWeights[len(a.clientModelWeights)-1])
}
