package entities

import (
	"github.com/Pro7ech/lattigo/rlwe"
	"math"
	"time"
)

type ClientModel struct {
	ClientUrl       string    `json:"client_url"`
	ModelName       string    `json:"model_name"`
	Length          int       `json:"model_shape"`
	Weights         [][]byte  `json:"weights"`
	LastModelUpdate time.Time `json:"-"`
}

func (m *ClientModel) GetIdentifier() string {
	return m.ClientUrl + m.ModelName
}

func (m *ClientModel) WeightsAsCiphertext() []*rlwe.Ciphertext {
	ciphertexts := make([]*rlwe.Ciphertext, len(m.Weights))
	for i, w := range m.Weights {
		ciphertext := rlwe.Ciphertext{}
		ciphertext.UnmarshalBinary(w)
		ciphertexts[i] = &ciphertext
	}
	return ciphertexts
}

func (m *ClientModel) SetNewWeights(weights [][]byte, length int) {
	m.Weights = weights
	m.Length = length
	m.LastModelUpdate = time.Now()
}

func (m *ClientModel) SetCiphertextWeights(weights []*rlwe.Ciphertext) {
	m.Weights = make([][]byte, len(weights))
	for i, w := range weights {
		binary, err := w.MarshalBinary()
		if err != nil {
			return
		}
		m.Weights[i] = binary
	}
}

type ModelClientUpdate struct {
	ModelName string    `json:"model_name"`
	Length    int       `json:"length"`
	Weights   []float32 `json:"weights"`
}

func (m *ModelClientUpdate) AddWeights(weights []float64) {
	m.Length = len(weights)
	m.Weights = make([]float32, len(weights))
	for i, w := range weights {
		m.Weights[i] = float32(w)
		if m.Weights[i] == float32(math.Inf(+1)) {
			m.Weights[i] = math.MaxFloat32
		}
		if m.Weights[i] == float32(math.Inf(-1)) {
			m.Weights[i] = -math.MaxFloat32
		}
	}
}
