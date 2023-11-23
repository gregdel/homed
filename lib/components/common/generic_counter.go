package common

import (
	"github.com/gregdel/homed/lib/components"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	components.Register(components.TypeGenericCounter, NewGenericCounter)
}

// NewGenericCounter returns a new generic sensor
func NewGenericCounter() components.Component {
	return &GenericCounter{}
}

// GenericCounter represents a generic sensor
type GenericCounter struct {
	GenericSensor
}

// Type implements the Component interface
func (g *GenericCounter) Type() components.Type {
	return components.TypeGenericCounter
}

// Collectors implements the Component interface
func (g *GenericCounter) Collectors(labels prometheus.Labels) []prometheus.Collector {
	return []prometheus.Collector{
		components.CounterCollector("generic_counter", labels,
			func() float64 { return g.SensorValue() },
		),
	}
}
