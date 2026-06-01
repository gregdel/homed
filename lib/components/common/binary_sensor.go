package common

import (
	"fmt"
	"sync"

	"github.com/gregdel/homed/lib/components"
	"github.com/prometheus/client_golang/prometheus"
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

	mu sync.RWMutex
	on bool
}

type BinarySensorSnapshot struct {
	SnapshotBase
	On bool `json:"on"`
}

// Type implements the Component interface
func (b *BinarySensor) Type() components.Type {
	return components.TypeBinarySensor
}

// IsOn implements the BinarySensor interface
func (b *BinarySensor) IsOn() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return b.on
}

func (b *BinarySensor) SetOn(on bool) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.on = on
}

// Update implements the Component interface
func (b *BinarySensor) Update(value []byte) error {
	switch string(value) {
	case "ON":
		b.SetOn(true)
	case "OFF":
		b.SetOn(false)
	default:
		return fmt.Errorf("binary_sensor: invalid payload: %s", value)
	}

	return nil
}

func (b *BinarySensor) ValuesSnapshot() any {
	return BinarySensorSnapshot{
		SnapshotBase: b.SnapshotBase(),
		On:           b.IsOn(),
	}
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
