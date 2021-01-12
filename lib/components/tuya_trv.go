package components

import (
	"encoding/json"

	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	register(TypeTuyaTRV, NewTuyaTRV)
}

// TuyaTRV is a component that handles temperatures
type TuyaTRV struct {
	baseComponent

	HeatingSetpoint float64 `json:"current_heating_setpoint"`
	Temperature     float64 `json:"local_temperature"`
	Position        float64 `json:"position"`
	BatteryLow      bool    `json:"battery_low"`
}

// NewTuyaTRV returns a new component for Tuya TRVs
func NewTuyaTRV() Component {
	return &TuyaTRV{}
}

// Type implements the Component interface
func (s *TuyaTRV) Type() Type {
	return TypeTuyaTRV
}

// Collectors implements the Component interface
func (s *TuyaTRV) Collectors(labels prometheus.Labels) []prometheus.Collector {
	prefix := "homed_tuya_trv"
	return []prometheus.Collector{
		prometheus.NewGaugeFunc(
			prometheus.GaugeOpts{
				Name:        prefix + "temperature",
				ConstLabels: labels,
			},
			func() float64 { return s.Temperature },
		),
		prometheus.NewGaugeFunc(
			prometheus.GaugeOpts{
				Name:        prefix + "heating_set_point",
				ConstLabels: labels,
			},
			func() float64 { return s.HeatingSetpoint },
		),
		prometheus.NewGaugeFunc(
			prometheus.GaugeOpts{
				Name:        prefix + "position",
				ConstLabels: labels,
			},
			func() float64 { return s.Position },
		),
		prometheus.NewGaugeFunc(
			prometheus.GaugeOpts{
				Name:        prefix + "battery_low",
				ConstLabels: labels,
			},
			func() float64 {
				if s.BatteryLow {
					return 1
				}
				return 0
			},
		),
	}
}

// Update implements the Component interface
func (s *TuyaTRV) Update(value []byte) error {
	return json.Unmarshal(value, s)
}
