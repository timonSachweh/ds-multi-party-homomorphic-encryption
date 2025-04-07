package entities

import (
	"time"

	"github.com/tuneinsight/lattigo/v6/core/rlwe"
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
