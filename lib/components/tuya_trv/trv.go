package trv

import (
	"encoding/json"
	"fmt"

	"github.com/gregdel/homed/lib/components"
	base "github.com/gregdel/homed/lib/components/base_component"
	"github.com/prometheus/client_golang/prometheus"
)

// Make sure that the module is a temperature controller
var _ components.TemperatureController = (*TuyaTRV)(nil)

func init() {
	components.Register(components.TypeTuyaTRV, New)
}

// TuyaTRV is a component that handles temperatures
type TuyaTRV struct {
	base.Component

	HeatingSetpoint  float64 `json:"current_heating_setpoint"`
	LocalTemperature float64 `json:"local_temperature"`
	Position         float64 `json:"position"`
	BatteryLow       bool    `json:"battery_low"`
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
			func() float64 { return t.LocalTemperature },
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

// SetTemperature implements the TemperatureGetterSetter interface
func (t *TuyaTRV) SetTemperature(float64) error {
	return components.ErrNotImplemented
}

// Temperature implements the TemperatureGetter interface
func (t *TuyaTRV) Temperature() (float64, error) {
	return t.LocalTemperature, nil
}

// TemperatureTarget implements the TemperatureController interface
func (t *TuyaTRV) TemperatureTarget() (float64, error) {
	return t.HeatingSetpoint, nil
}

// SetTemperatureTarget implements the TemperatureController interface
func (t *TuyaTRV) SetTemperatureTarget(temperature float64) error {
	if temperature == t.HeatingSetpoint {
		return nil
	}

	t.HeatingSetpoint = temperature

	data := fmt.Sprintf("%.2f", temperature)
	return t.WriteCommand([]byte(data))
}

// TemperatureMode implements the TemperatureController interface
func (t *TuyaTRV) TemperatureMode() (components.TemperatureMode, error) {
	return components.TemperatureModeAuto, components.ErrNotImplemented
}

// SetTemperatureMode implements the TemperatureController interface
func (t *TuyaTRV) SetTemperatureMode(components.TemperatureMode) error {
	return components.ErrNotImplemented
}
