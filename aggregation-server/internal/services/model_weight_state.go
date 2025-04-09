package services

import (
	"errors"
	"github.com/tuneinsight/lattigo/v6/core/rlwe"
	"maps"
	"slices"
	"time"

	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/aggregationserver/internal/entities"
)

type ModelWeightStateManager interface {
	AddModelWeightsFromClient(entities.ClientModel)
	AddClient(client entities.ClientModel) error
	ModelCanAggregate() bool
	GetClients() ([]entities.ClientModel, [][]*rlwe.Ciphertext, []string, int)
	GetIdentifiers() []string
	GetClientUrls() []string
	ResetClientWeights()
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
	m.clients[modelWeights.GetIdentifier()] = model
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

func (m *ModelWeightStateImpl) ResetClientWeights() {
	for key, client := range m.clients {
		client.Weights = nil
		client.LastModelUpdate = time.Time{}
		m.clients[key] = client
	}
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

func (m *ModelWeightStateImpl) GetClients() ([]entities.ClientModel, [][]*rlwe.Ciphertext, []string, int) {
	clients := make([]entities.ClientModel, 0)
	modelWeights := make([][]*rlwe.Ciphertext, 0)
	clientUrls := make([]string, 0)
	weightLength := 0
	for _, client := range m.clients {
		clients = append(clients, client)
		modelWeights = append(modelWeights, client.WeightsAsCiphertext())
		clientUrls = append(clientUrls, client.ClientUrl)
		if weightLength == 0 {
			weightLength = client.Length
		}
	}
	return clients, modelWeights, clientUrls, weightLength
}

func (m *ModelWeightStateImpl) ModelCanAggregate() bool {
	if len(m.clients) < m.minRequiredClients {
		return false
	}

	for _, modelWeights := range m.clients {
		if modelWeights.LastModelUpdate.IsZero() {
			return false
		}
	}

	return true
}
