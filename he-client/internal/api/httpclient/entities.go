package httpclient

import "github.com/tuneinsight/lattigo/v6/core/rlwe"

type MLModelWeights struct {
	ModelName string          `json:"model_name"`
	Length    int             `json:"model_shape"`
	Weights   rlwe.Ciphertext `json:"weights"`
}
