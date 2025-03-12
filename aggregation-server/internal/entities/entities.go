package entities

import "github.com/tuneinsight/lattigo/v6/core/rlwe"

type MLModelWeights struct {
	ModelName string   `json:"model_name"`
	Length    int      `json:"model_shape"`
	Weights   [][]byte `json:"weights"`
}

func (m *MLModelWeights) WeightsAsCiphertext() []*rlwe.Ciphertext {
	ciphertexts := make([]*rlwe.Ciphertext, len(m.Weights))
	for i, w := range m.Weights {
		ciphertext := rlwe.Ciphertext{}
		ciphertext.UnmarshalBinary(w)
		ciphertexts[i] = &ciphertext
	}
	return ciphertexts
}
