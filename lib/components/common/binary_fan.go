package common

import (
	"github.com/gregdel/homed/lib/components"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	components.Register(components.TypeBinaryFan, NewBinaryFan)
}

// BinaryFan represents a generic binary fan actionnable as switch
type BinaryFan struct {
	Switch
}

// NewBinaryFan returns a new binary fan
func NewBinaryFan() components.Component {
	return &BinaryFan{}
}

// Type implements the Component interface
func (b *BinaryFan) Type() components.Type {
	return components.TypeBinaryFan
}

// Collectors implements the Component interface
func (b *BinaryFan) Collectors(labels prometheus.Labels) []prometheus.Collector {
	return []prometheus.Collector{
		components.GaugeCollector("binary_fan", labels,
			func() float64 {
				if b.IsOn() {
					return 1
				}
				return 0
			},
		),
	}
}
