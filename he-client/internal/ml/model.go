package ml

import (
	"math/rand"
	"time"
)

type Model interface {
	Train()
	Predict()
	AsFloatVector() []float64
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

}

func (m *ModelImpl) Predict() {

}

func (m *ModelImpl) AsFloatVector() []float64 {
	return m.module
}

func generateRandomMatrix(length int, min, max float64) []float64 {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	matrix := make([]float64, length)
	for i := range matrix {
		matrix[i] = min + rand.Float64()*(max-min)
	}
	return matrix
}
