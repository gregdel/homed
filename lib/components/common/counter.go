package common

import (
	"github.com/gregdel/homed/lib/components"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	components.Register(components.TypeCounter, NewCounter)
}

// Counter is a simple coutner used to count stuff, it should not be displayed
// in the UI for now
type Counter struct {
	GenericSensor
}

// NewCounter returns a new power meter
func NewCounter() components.Component {
	return &Counter{}
}

// Type implements the Component interface
func (c *Counter) Type() components.Type {
	return components.TypeCounter
}

// Collectors implements the Component interface
func (c *Counter) Collectors(labels prometheus.Labels) []prometheus.Collector {
	return []prometheus.Collector{
		prometheus.NewCounterFunc(
			prometheus.CounterOpts{
				Name:        "homed_counter",
				ConstLabels: labels,
			},
			func() float64 { return c.Value },
		),
	}
}
