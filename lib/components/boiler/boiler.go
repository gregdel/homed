package boiler

import (
	"fmt"
	"time"

	"github.com/gregdel/homed/lib/components"
	base "github.com/gregdel/homed/lib/components/base_component"
	"github.com/prometheus/client_golang/prometheus"
)

const boilerCooldownDuration = 5 * time.Minute

func init() {
	components.Register("boiler", NewBoiler)
}

// Boiler is a component that controls the boiler
type Boiler struct {
	base.Component

	On              bool       `json:"on"`
	LastStateChange *time.Time `json:"last_state_change"`
}

// NewBoiler returns a new status component
func NewBoiler() components.Component {
	return &Boiler{}
}

// Type implements the Component interface
func (b *Boiler) Type() components.Type {
	return components.TypeBoiler
}

// Collectors implements the Component interface
func (b *Boiler) Collectors(labels prometheus.Labels) []prometheus.Collector {
	return []prometheus.Collector{
		prometheus.NewGaugeFunc(
			prometheus.GaugeOpts{
				Name:        "homed_boiler",
				ConstLabels: labels,
			},
			func() float64 {
				if b.On {
					return 1
				}
				return 0
			},
		),
	}
}

func (b *Boiler) stateFromData(value []byte) bool {
	data := string(value)
	if data == "ON" {
		return true
	}
	return false
}

// WriteCommand implements the Component interface
func (b *Boiler) WriteCommand(data []byte) error {
	newState := b.stateFromData(data)
	if newState == b.On {
		return nil
	}

	now := time.Now()
	if b.LastStateChange == nil {
		b.LastStateChange = &now
	} else {
		if b.LastStateChange.Add(boilerCooldownDuration).Before(now) {
			return fmt.Errorf("components: boiler: last change is to recent")
		}
	}

	return b.Component.WriteCommand(data)
}

// Update implements the Component interface
func (b *Boiler) Update(value []byte) error {
	b.On = b.stateFromData(value)
	return nil
}
