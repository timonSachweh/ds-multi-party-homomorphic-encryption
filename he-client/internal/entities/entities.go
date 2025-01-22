package entities

import "github.com/tuneinsight/lattigo/v6/core/rlwe"

type MLModelWeights struct {
	ModelName string `json:"model_name"`
	Length    int    `json:"model_shape"`
	Weights   []byte `json:"weights"`
}

func (m *MLModelWeights) WeightsAsCiphertext() *rlwe.Ciphertext {
	var ciphertext rlwe.Ciphertext
	ciphertext.UnmarshalBinary(m.Weights)
	return &ciphertext
}
