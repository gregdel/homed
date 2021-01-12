package components

import (
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	register(TypeHumidity, NewHumidity)
}

// Humidity is a component that handles temperatures
type Humidity struct {
	baseComponent
	float64Component
}

// NewHumidity returns a new humidity component
func NewHumidity() Component {
	return &Humidity{}
}

// Type implements the Component interface
func (s *Humidity) Type() Type {
	return TypeHumidity
}

// Collectors implements the Component interface
func (s *Humidity) Collectors(labels prometheus.Labels) []prometheus.Collector {
	return singleCollector(s, labels, func() float64 { return s.Value })
}
