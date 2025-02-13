package ml

import (
	"log"
	"math/rand"
	"time"
)

type Model interface {
	Train()
	Predict()
	Load(path string)
	AsFloatVector() []float64
	UpdateWeights(weights []float64)
}

type ModelImpl struct {
	module []float64
}

func NewModel() Model {
	m := &ModelImpl{}
	m.module = generateRandomMatrix(5, -1, 4)
	return m
}

func (m *ModelImpl) Train() {
	log.Println("Training model")
}

func (m *ModelImpl) Predict() {

}

func (m *ModelImpl) Load(path string) {

}

func (m *ModelImpl) AsFloatVector() []float64 {
	return m.module
}

func (m *ModelImpl) UpdateWeights(weights []float64) {
	m.module = weights
	log.Println(m.module)
}

func generateRandomMatrix(length int, min, max float64) []float64 {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	matrix := make([]float64, length)
	for i := range matrix {
		matrix[i] = min + rand.Float64()*(max-min)
	}
	return matrix
}
