package aggregationupdate

// Service defines the interface for aggregation update operations.
type AggregationUpdateService interface {
}

type aggregationUpdateServiceImpl struct {
}

// NewAggregationUpdateService creates a new instance of aggregationUpdateServiceImpl.
func NewAggregationUpdateService() AggregationUpdateService {
	return &aggregationUpdateServiceImpl{}
}
