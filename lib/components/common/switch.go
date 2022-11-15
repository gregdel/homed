package common

import (
	"github.com/gregdel/homed/lib/components"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	components.Register(components.TypeSwitch, NewSwitch)
}

// Switch represents a generic switch
type Switch struct {
	BinarySensor
}

// NewSwitch returns a new switch
func NewSwitch() components.Component {
	return &Switch{}
}

// TurnOn implements the switch interface
func (s *Switch) TurnOn() error {
	return s.WriteCommand([]byte("ON"))
}

// TurnOff implements the switch interface
func (s *Switch) TurnOff() error {
	return s.WriteCommand([]byte("OFF"))
}

// Toggle implements the Switch interface
func (s *Switch) Toggle() error {
	if s.IsOn() {
		return s.TurnOff()
	}

	return s.TurnOn()
}

// Type implements the Component interface
func (s *Switch) Type() components.Type {
	return components.TypeSwitch
}

// Collectors implements the Component interface
func (s *Switch) Collectors(labels prometheus.Labels) []prometheus.Collector {
	return []prometheus.Collector{
		components.GaugeCollector("switch", labels,
			func() float64 {
				if s.On {
					return 1
				}
				return 0
			},
		),
	}
}
