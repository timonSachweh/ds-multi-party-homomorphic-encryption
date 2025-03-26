package services

import (
	"errors"
	"maps"
	"slices"
	"time"

	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/aggregationserver/internal/entities"
)

type ModelWeightStateManager interface {
	AddModelWeightsFromClient(entities.MLModelWeights)
	GetAggregatedModelWeights() (entities.MLModelWeights, []string, error)
	GetIdentifiers() []string
}

type ModelWeightStateImpl struct {
	weights             map[string]entities.MLModelWeights
	updatedModelWeights []entities.MLModelWeights
	minRequiredClients  int
}

func NewModelWeightStateManager(minRequiredClients int) ModelWeightStateManager {
	return &ModelWeightStateImpl{
		weights:             make(map[string]entities.MLModelWeights),
		updatedModelWeights: make([]entities.MLModelWeights, 0),
		minRequiredClients:  minRequiredClients,
	}
}

func (m *ModelWeightStateImpl) AddModelWeightsFromClient(modelWeights entities.MLModelWeights) {
	m.weights[modelWeights.GetIdentifier()] = modelWeights
}

func (m *ModelWeightStateImpl) GetUniqueModelNames() []string {
	uniqueModelNames := make([]string, 0)
	for _, modelWeights := range m.weights {
		if slices.Contains(uniqueModelNames, modelWeights.ModelName) {
			continue
		}
		uniqueModelNames = append(uniqueModelNames, modelWeights.ModelName)
	}
	return uniqueModelNames
}

func (m *ModelWeightStateImpl) GetIdentifiers() []string {
	return slices.Collect(maps.Keys(m.weights))
}

func (m *ModelWeightStateImpl) GetAggregatedModelWeights() (entities.MLModelWeights, []string, error) {
	if len(m.weights) < m.minRequiredClients {
		return entities.MLModelWeights{}, nil, errors.New("Not enough clients to aggregate")
	}

	aggregatedResult := entities.MLModelWeights{}
	clientUrls := make([]string, 0)
	for _, modelWeights := range m.weights {
		aggregatedResult = modelWeights

		modelWeights.LastModelUpdate = time.Now()
		m.updatedModelWeights = append(m.updatedModelWeights, modelWeights)

		clientUrls = append(clientUrls, modelWeights.ClientUrl)
	}
	m.weights = make(map[string]entities.MLModelWeights)

	return aggregatedResult, clientUrls, nil
}
