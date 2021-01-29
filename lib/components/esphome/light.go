package esphome

import (
	"encoding/json"

	"github.com/gregdel/homed/lib/components"
	"github.com/gregdel/homed/lib/components/common"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	components.Register(components.TypeESPHomeLight, NewLight)
}

// Light is a component that controls a light
type Light struct {
	common.Component
	common.Switch
}

// NewLight returns a new light component
func NewLight() components.Component {
	payloadOn, _ := json.Marshal(lightState{State: "ON"})
	payloadOff, _ := json.Marshal(lightState{State: "OFF"})

	return &Light{
		Switch: common.NewSwitch(payloadOn, payloadOff),
	}
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

// Update implements the Component interface
func (l *Light) Update(value []byte) error {
	ls := lightState{}
	if err := json.Unmarshal(value, &ls); err != nil {
		return err
	}

	l.On = ls.isOn()
	return nil
}
