package common

import (
	"github.com/gregdel/homed/lib/components"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	components.Register(components.TypePowerMeter, NewPowerMeter)
}

// PowerMeter is a sensor that reads power
type PowerMeter struct {
	GenericSensor
}

// NewPowerMeter returns a new power meter
func NewPowerMeter() components.Component {
	return &PowerMeter{}
}

// Type implements the Component interface
func (p *PowerMeter) Type() components.Type {
	return components.TypePowerMeter
}

// Collectors implements the Component interface
func (p *PowerMeter) Collectors(labels prometheus.Labels) []prometheus.Collector {
	return []prometheus.Collector{
		components.GaugeCollector("power", labels,
			func() float64 { return p.SensorValue() },
		),
	}
}
