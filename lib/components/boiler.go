package components

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/prometheus/client_golang/prometheus"
)

const boilerCooldownDuration = 5 * time.Minute

func init() {
	register(TypeBoiler, NewBoiler)
}

// Boiler is a component that controls the boiler
type Boiler struct {
	baseComponent

	On              bool       `json:"on"`
	LastStateChange *time.Time `json:"last_state_change"`
}

// NewBoiler returns a new status component
func NewBoiler() Component {
	return &Boiler{}
}

// Type implements the Component interface
func (s *Boiler) Type() Type {
	return TypeBoiler
}

// Collectors implements the Component interface
func (s *Boiler) Collectors(labels prometheus.Labels) []prometheus.Collector {
	return []prometheus.Collector{
		prometheus.NewGaugeFunc(
			prometheus.GaugeOpts{
				Name:        "homed_boiler",
				ConstLabels: labels,
			},
			func() float64 {
				if s.On {
					return 1
				}
				return 0
			},
		),
	}
}

func (s *Boiler) stateFromData(value []byte) bool {
	data := string(value)
	if data == "ON" {
		return true
	}
	return false
}

// WriteCommand implements the Component interface
func (s *Boiler) WriteCommand(client mqtt.Client, data []byte) error {
	newState := s.stateFromData(data)
	if newState == s.On {
		return nil
	}

	now := time.Now()
	if s.LastStateChange == nil {
		s.LastStateChange = &now
	} else {
		if s.LastStateChange.Add(boilerCooldownDuration).Before(now) {
			return fmt.Errorf("components: boiler: last change is to recent")
		}
	}

	return s.baseComponent.WriteCommand(client, data)
}

// Update implements the Component interface
func (s *Boiler) Update(value []byte) error {
	s.On = s.stateFromData(value)
	return nil
}
