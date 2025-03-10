package entities

import "github.com/tuneinsight/lattigo/v6/core/rlwe"

type MLModelWeights struct {
	ModelName string    `json:"model_name"`
	Length    int       `json:"model_shape"`
	Weights   []float32 `json:"weights"`
}

func (m *MLModelWeights) ToDataSpaceModelWeights() *DataSpaceModelWeights {
	return &DataSpaceModelWeights{
		ModelName: m.ModelName,
		Length:    m.Length,
	}
}

type DataSpaceModelWeights struct {
	ModelName string `json:"model_name"`
	Length    int    `json:"model_shape"`
	Weights   []byte `json:"weights"`
}

func (m *DataSpaceModelWeights) WeightsAsCiphertext() *rlwe.Ciphertext {
	var ciphertext rlwe.Ciphertext
	ciphertext.UnmarshalBinary(m.Weights)
	return &ciphertext
}

type PredictionRequest struct {
	Data [][]float32 `json:"data"`
}

type PredictionResponse struct {
	Prediction []float32 `json:"prediction"`
}
