package services

import (
	"errors"
	"maps"
	"slices"
	"time"

	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/aggregationserver/internal/entities"
)

type ModelWeightStateManager interface {
	AddModelWeightsFromClient(entities.ClientModel)
	AddClient(client entities.ClientModel) error
	GetAggregatedModelWeights() (entities.ClientModel, []string, error)
	GetIdentifiers() []string
	GetClientUrls() []string
}

type ModelWeightStateImpl struct {
	clients             map[string]entities.ClientModel
	updatedModelWeights []entities.ClientModel
	minRequiredClients  int
}

func NewModelWeightStateManager(minRequiredClients int) ModelWeightStateManager {
	return &ModelWeightStateImpl{
		clients:             make(map[string]entities.ClientModel),
		updatedModelWeights: make([]entities.ClientModel, 0),
		minRequiredClients:  minRequiredClients,
	}
}

func (m *ModelWeightStateImpl) AddModelWeightsFromClient(modelWeights entities.ClientModel) {
	model := m.clients[modelWeights.GetIdentifier()]
	model.SetNewWeights(modelWeights.Weights, modelWeights.Length)
}

func (m *ModelWeightStateImpl) AddClient(client entities.ClientModel) error {
	if client.GetIdentifier() == "" {
		return errors.New("client identifier is empty")
	}
	if _, ok := m.clients[client.GetIdentifier()]; ok {
		return errors.New("client already exists")
	}
	m.clients[client.GetIdentifier()] = client
	return nil
}

func (m *ModelWeightStateImpl) GetClientUrls() []string {
	clientUrls := make([]string, 0)
	for _, modelWeights := range m.clients {
		clientUrls = append(clientUrls, modelWeights.ClientUrl)
	}
	return clientUrls
}

func (m *ModelWeightStateImpl) GetUniqueModelNames() []string {
	uniqueModelNames := make([]string, 0)
	for _, modelWeights := range m.clients {
		if slices.Contains(uniqueModelNames, modelWeights.ModelName) {
			continue
		}
		uniqueModelNames = append(uniqueModelNames, modelWeights.ModelName)
	}
	return uniqueModelNames
}

func (m *ModelWeightStateImpl) GetIdentifiers() []string {
	return slices.Collect(maps.Keys(m.clients))
}

func (m *ModelWeightStateImpl) GetAggregatedModelWeights() (entities.ClientModel, []string, error) {
	if len(m.clients) < m.minRequiredClients {
		return entities.ClientModel{}, nil, errors.New("not enough clients to aggregate")
	}

	aggregatedResult := entities.ClientModel{}
	clientUrls := make([]string, 0)
	for _, modelWeights := range m.clients {
		aggregatedResult = modelWeights

		modelWeights.LastModelUpdate = time.Now()
		m.updatedModelWeights = append(m.updatedModelWeights, modelWeights)

		clientUrls = append(clientUrls, modelWeights.ClientUrl)
	}
	m.clients = make(map[string]entities.ClientModel)

	return aggregatedResult, clientUrls, nil
}
