package components

import (
	"encoding/json"

	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	register(TypeTasmotaSwitch, NewTasmotaSwitch)
}

// TasmotaSwitch is a component that controls the boiler
type TasmotaSwitch struct {
	baseComponent

	On bool `json:"on"`
}

// NewTasmotaSwitch returns a new status component
func NewTasmotaSwitch() Component {
	return &TasmotaSwitch{}
}

// Type implements the Component interface
func (s *TasmotaSwitch) Type() Type {
	return TypeTasmotaSwitch
}

// Collectors implements the Component interface
func (s *TasmotaSwitch) Collectors(labels prometheus.Labels) []prometheus.Collector {
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
func (s *TasmotaSwitch) Update(value []byte) error {
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
