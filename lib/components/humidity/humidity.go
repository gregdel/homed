package humidity

import (
	"github.com/gregdel/homed/lib/components"
	"github.com/gregdel/homed/lib/components/common"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	components.Register(components.TypeHumidity, New)
}

// Humidity is a component that handles temperatures
type Humidity struct {
	common.Component
	common.GenericSensor
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
