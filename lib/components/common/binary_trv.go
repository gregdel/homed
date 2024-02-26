package common

import (
	"fmt"
	"strconv"

	"github.com/gregdel/homed/lib/components"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/atomic"
)

func init() {
	components.Register(components.TypeBinaryTRV, NewBinaryTRV)
}

// BinaryTRV represents a generic TRV that can be turned on or off
type BinaryTRV struct {
	Switch

	currentSetPoint atomic.Float64
	OnTemperature   int
	OffTemperature  int
}

// NewBinaryTRV returns a new binary light
func NewBinaryTRV() components.Component {
	return &BinaryTRV{
		// For now those are constants, we can add them later in the config
		OnTemperature:  28,
		OffTemperature: 8,
	}
}

// Update implements the Component interface
func (b *BinaryTRV) Update(value []byte) error {
	v, err := strconv.ParseFloat(string(value), 64)
	if err != nil {
		return fmt.Errorf("binary_trv: invalid payload: %s", value)
	}

	b.On.Store(v == float64(b.OnTemperature))
	b.currentSetPoint.Store(v)
	return nil
}

// TurnOn implements the switch interface
func (b *BinaryTRV) TurnOn() error {
	if b.IsOn() {
		return nil
	}

	return b.WriteCommand([]byte(strconv.Itoa(b.OnTemperature)))
}

// TurnOff implements the switch interface
func (b *BinaryTRV) TurnOff() error {
	if b.currentSetPoint.Load() == float64(b.OffTemperature) {
		return nil
	}

	return b.WriteCommand([]byte(strconv.Itoa(b.OffTemperature)))
}

// Type implements the Component interface
func (b *BinaryTRV) Type() components.Type {
	return components.TypeBinaryTRV
}

// Collectors implements the Component interface
func (b *BinaryTRV) Collectors(labels prometheus.Labels) []prometheus.Collector {
	return []prometheus.Collector{
		components.GaugeCollector("binary_trv", labels,
			func() float64 {
				if b.IsOn() {
					return 1
				}
				return 0
			},
		),
	}
}
