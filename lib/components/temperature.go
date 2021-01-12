package components

import (
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	register(TypeTemperature, NewTemperature)
}

// Temperature is a component that handles temperatures
type Temperature struct {
	baseComponent
	float64Component
}

// NewTemperature returns a new temperature component
func NewTemperature() Component {
	return &Temperature{}
}

// Type implements the Component interface
func (s *Temperature) Type() Type {
	return TypeTemperature
}

// Collectors implements the Component interface
func (s *Temperature) Collectors(labels prometheus.Labels) []prometheus.Collector {
	return singleCollector(s, labels, func() float64 { return s.Value })
}
