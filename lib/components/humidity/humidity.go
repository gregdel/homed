package humidity

import (
	"github.com/gregdel/homed/lib/components"
	base "github.com/gregdel/homed/lib/components/base_component"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	components.Register(components.TypeHumidity, New)
}

// Humidity is a component that handles temperatures
type Humidity struct {
	base.Component
	base.Float64Component
}

// New returns a new humidity component
func New() components.Component {
	return &Humidity{}
}

// Type implements the Component interface
func (h *Humidity) Type() components.Type {
	return components.TypeHumidity
}

// Collectors implements the Component interface
func (h *Humidity) Collectors(labels prometheus.Labels) []prometheus.Collector {
	return components.SingleCollector(h, labels, func() float64 { return h.Value })
}
