package aggregation

type AggregationService interface {
}

type AggregationServiceImpl struct {
}

// NewAggregationService creates a new instance of AggregationServiceImpl.
func NewAggregationService() AggregationService {
	return &AggregationServiceImpl{}
}
