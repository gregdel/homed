package common

import (
	"fmt"

	"github.com/gregdel/homed/lib/components"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/atomic"
)

func init() {
	components.Register(components.TypeBinarySensor, NewBinarySensor)
}

// NewBinarySensor returns a new binary sensor
func NewBinarySensor() components.Component {
	return &BinarySensor{}
}

// BinarySensor represents a generic sensor
type BinarySensor struct {
	Component

	On atomic.Bool `json:"on"`
}

// Type implements the Component interface
func (b *BinarySensor) Type() components.Type {
	return components.TypeBinarySensor
}

// IsOn implements the BinarySensor interface
func (b *BinarySensor) IsOn() bool {
	return b.On.Load()
}

// Update implements the Component interface
func (b *BinarySensor) Update(value []byte) error {
	switch string(value) {
	case "ON":
		b.On.Store(true)
	case "OFF":
		b.On.Store(false)
	default:
		return fmt.Errorf("binary_sensor: invalid payload: %s", value)
	}

	return nil
}

// Collectors implements the Component interface
func (b *BinarySensor) Collectors(labels prometheus.Labels) []prometheus.Collector {
	return []prometheus.Collector{
		components.GaugeCollector("binary_sensor", labels,
			func() float64 {
				if b.IsOn() {
					return 1
				}

				return 0
			},
		),
	}
}
