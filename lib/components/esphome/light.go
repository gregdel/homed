package esphome

import (
	"encoding/json"

	"github.com/gregdel/homed/lib/components"
	base "github.com/gregdel/homed/lib/components/base_component"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	components.Register(components.TypeESPHomeLight, NewLight)
}

// Light is a component that controls a light
type Light struct {
	base.Component
	On bool `json:"on"`
}

// NewLight returns a new light component
func NewLight() components.Component {
	return &Light{}
}

// Type implements the Component interface
func (l *Light) Type() components.Type {
	return components.TypeESPHomeLight
}

// Collectors implements the Component interface
func (l *Light) Collectors(labels prometheus.Labels) []prometheus.Collector {
	return []prometheus.Collector{
		prometheus.NewGaugeFunc(
			prometheus.GaugeOpts{
				Name:        "homed_light",
				ConstLabels: labels,
			},
			func() float64 {
				if l.On {
					return 1
				}
				return 0
			},
		),
	}
}

type lightState struct {
	State string `json:"state"`
}

func (ls *lightState) isOn() bool {
	if ls.State == "ON" {
		return true
	}

	return false
}

func (l *Light) payload(on bool) ([]byte, error) {
	state := "OFF"
	if on {
		state = "ON"
	}
	return json.Marshal(lightState{State: state})
}

// Update implements the Component interface
func (l *Light) Update(value []byte) error {
	ls := lightState{}
	if err := json.Unmarshal(value, &ls); err != nil {
		return err
	}

	l.On = ls.isOn()
	return nil
}

// SetOn implements the Switch interface
func (l *Light) SetOn() error {
	payload, err := l.payload(true)
	if err != nil {
		return err
	}
	return l.WriteCommand(payload)
}

// SetOff implements the Switch interface
func (l *Light) SetOff() error {
	payload, err := l.payload(false)
	if err != nil {
		return err
	}
	return l.WriteCommand(payload)
}

// Set implements the Switch interface
func (l *Light) Set(state bool) error {
	if state {
		return l.SetOn()
	}

	return l.SetOff()
}

// Toggle implements the Switch interface
func (l *Light) Toggle() error {
	if l.On {
		return l.SetOn()
	}

	return l.SetOff()
}

// IsOn implements the Switch interface
func (l *Light) IsOn() bool {
	return l.On
}
