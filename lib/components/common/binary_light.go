package common

import (
	"github.com/gregdel/homed/lib/components"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	components.Register(components.TypeBinaryLight, NewBinaryLight)
}

// BinaryLight represents a generic binary light actionnable with a switch
type BinaryLight struct {
	Switch
}

// NewBinaryLight returns a new binary light
func NewBinaryLight() components.Component {
	return &BinaryLight{}
}

// Type implements the Component interface
func (b *BinaryLight) Type() components.Type {
	return components.TypeBinaryLight
}

// Collectors implements the Component interface
func (b *BinaryLight) Collectors(labels prometheus.Labels) []prometheus.Collector {
	return []prometheus.Collector{
		components.GaugeCollector("binary_light", labels,
			func() float64 {
				if b.IsOn() {
					return 1
				}
				return 0
			},
		),
	}
}
