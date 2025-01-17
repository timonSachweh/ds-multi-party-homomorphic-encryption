package ml

type MLService interface {
	Train()
	Predict()
	GetModel() Model
}

type MLServiceImpl struct {
	model Model
}

func NewMLService() MLService {
	return &MLServiceImpl{
		model: NewModel(),
	}
}

func (m *MLServiceImpl) Train() {
	m.model.Train()
}

func (m *MLServiceImpl) Predict() {
	m.model.Predict()
}

func (m *MLServiceImpl) GetModel() Model {
	return m.model
}
