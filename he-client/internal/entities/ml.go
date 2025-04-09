package entities

import "github.com/tuneinsight/lattigo/v6/core/rlwe"

type MLModelWeights struct {
	ModelName string    `json:"model_name"`
	Length    int       `json:"model_shape"`
	Weights   []float32 `json:"weights"`
}

func (m *MLModelWeights) ToDataSpaceModelWeights(url string) *ClientModel {
	return &ClientModel{
		ClientUrl: url,
		ModelName: m.ModelName,
		Length:    m.Length,
	}
}

type ClientModel struct {
	ClientUrl string   `json:"client_url"`
	ModelName string   `json:"model_name"`
	Length    int      `json:"model_shape"`
	Weights   [][]byte `json:"weights"`
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

type ModelClientUpdate struct {
	ModelName string    `json:"model_name"`
	Length    int       `json:"length"`
	Weights   []float32 `json:"weights"`
}

type PredictionRequest struct {
	Data [][]float32 `json:"data"`
}

type PredictionResponse struct {
	Prediction []float32 `json:"prediction"`
}
