package entities

type MLModelWeights struct {
	ModelName string `json:"model_name"`
	Length    int    `json:"model_shape"`
	Weights   []byte `json:"weights"`
}
