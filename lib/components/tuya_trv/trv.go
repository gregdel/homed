package trv

import (
	"encoding/json"

	"github.com/gregdel/homed/lib/components"
	base "github.com/gregdel/homed/lib/components/base_component"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	components.Register(components.TypeTuyaTRV, New)
}

// TuyaTRV is a component that handles temperatures
type TuyaTRV struct {
	base.Component

	HeatingSetpoint float64 `json:"current_heating_setpoint"`
	Temperature     float64 `json:"local_temperature"`
	Position        float64 `json:"position"`
	BatteryLow      bool    `json:"battery_low"`
}

// New returns a new component for Tuya TRVs
func New() components.Component {
	return &TuyaTRV{}
}

// Type implements the Component interface
func (t *TuyaTRV) Type() components.Type {
	return components.TypeTuyaTRV
}

// Collectors implements the Component interface
func (t *TuyaTRV) Collectors(labels prometheus.Labels) []prometheus.Collector {
	prefix := "homed_tuya_trv"
	return []prometheus.Collector{
		prometheus.NewGaugeFunc(
			prometheus.GaugeOpts{
				Name:        prefix + "temperature",
				ConstLabels: labels,
			},
			func() float64 { return t.Temperature },
		),
		prometheus.NewGaugeFunc(
			prometheus.GaugeOpts{
				Name:        prefix + "heating_set_point",
				ConstLabels: labels,
			},
			func() float64 { return t.HeatingSetpoint },
		),
		prometheus.NewGaugeFunc(
			prometheus.GaugeOpts{
				Name:        prefix + "position",
				ConstLabels: labels,
			},
			func() float64 { return t.Position },
		),
		prometheus.NewGaugeFunc(
			prometheus.GaugeOpts{
				Name:        prefix + "battery_low",
				ConstLabels: labels,
			},
			func() float64 {
				if t.BatteryLow {
					return 1
				}
				return 0
			},
		),
	}
}

// Update implements the Component interface
func (t *TuyaTRV) Update(value []byte) error {
	return json.Unmarshal(value, t)
}
