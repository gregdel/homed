package tasmota

import (
	"encoding/json"

	"github.com/gregdel/homed/lib/components"
	"github.com/gregdel/homed/lib/components/common"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	components.Register(components.TypeTasmotaSwitch, New)
}

// Switch is a component that controls the boiler
type Switch struct {
	common.Component
	common.Switch
}

// New returns a new status component
func New() components.Component {
	return &Switch{
		Switch: common.NewSwitch([]byte("ON"), []byte("OFF")),
	}
}

// Type implements the Component interface
func (s *Switch) Type() components.Type {
	return components.TypeTasmotaSwitch
}

// Collectors implements the Component interface
func (s *Switch) Collectors(labels prometheus.Labels) []prometheus.Collector {
	return []prometheus.Collector{
		prometheus.NewGaugeFunc(
			prometheus.GaugeOpts{
				Name:        "homed_tasmota_switch",
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

// Update implements the Component interface
func (s *Switch) Update(value []byte) error {
	data := struct {
		Power string `json:"POWER"`
	}{}

	if err := json.Unmarshal(value, &data); err != nil {
		return err
	}

	if data.Power == "ON" {
		s.On = true
	} else {
		s.On = false
	}

	return nil
}
