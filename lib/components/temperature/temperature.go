package temperature

import (
	"github.com/gregdel/homed/lib/components"
	base "github.com/gregdel/homed/lib/components/base_component"
	"github.com/prometheus/client_golang/prometheus"
)

// Make sure that the module is a TemperatureGetter
var _ components.TemperatureGetter = (*Temperature)(nil)

func init() {
	components.Register(components.TypeTemperature, New)
}

// Temperature is a component that handles temperatures
type Temperature struct {
	base.Component
	base.GenericSensor
}

// New new temperature component
func New() components.Component {
	return &Temperature{}
}

// Type implements the Component interface
func (t *Temperature) Type() components.Type {
	return components.TypeTemperature
}

// Collectors implements the Component interface
func (t *Temperature) Collectors(labels prometheus.Labels) []prometheus.Collector {
	return components.SingleCollector(t, labels, func() float64 { return t.Value })
}

// Temperature implements the TemperatureGetter interface
func (t *Temperature) Temperature() (float64, error) {
	return t.Value, nil
}
