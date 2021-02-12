package trv

import (
	"encoding/json"
	"math"

	"github.com/gregdel/homed/lib/components"
	"github.com/gregdel/homed/lib/components/common"
	"github.com/prometheus/client_golang/prometheus"
)

// Make sure that the module is a temperature controller
var _ components.TemperatureController = (*TuyaTRV)(nil)

func init() {
	components.Register(components.TypeTuyaTRV, New)
}

// Force the trv to be fully open if we need to heat and we are far from the
// wanted temperature
const forceModeDiff = 0.5

// ForceMode represents the force modes
type ForceMode string

// SystemMode represents the system modes
type SystemMode string

// Force modes
var (
	ForceModeNormal ForceMode = "normal"
	ForceModeOpen   ForceMode = "open"
	ForceModeClose  ForceMode = "close"

	SystemModeAuto SystemMode = "auto"
	SystemModeHeat SystemMode = "heat"
	SystemModeOff  SystemMode = "off"
)

// TuyaTRV is a component that handles temperatures
type TuyaTRV struct {
	common.Component

	HeatingSetpoint             float64    `json:"current_heating_setpoint"`
	LocalTemperatureCalibration float64    `json:"local_temperature_calibration"`
	LocalTemperature            float64    `json:"local_temperature"`
	Position                    float64    `json:"position"`
	BatteryLow                  bool       `json:"battery_low"`
	Force                       ForceMode  `json:"force"`
	Mode                        SystemMode `json:"system_mode"`
}

// New returns a new component for Tuya TRVs
func New() components.Component {
	return &TuyaTRV{}
}

// Type implements the Component interface
func (t *TuyaTRV) Type() components.Type {
	return components.TypeTuyaTRV
}

// PostUpdate implements the Component interface
func (t *TuyaTRV) PostUpdate() error {
	if err := t.Component.PostUpdate(); err != nil {
		return err
	}

	return t.updateForceMode()
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

func (t *TuyaTRV) write(input interface{}) error {
	data, err := json.Marshal(input)
	if err != nil {
		return err
	}

	return t.WriteCommand(data)
}

// Format the float to 1 digit of precision
func formatFloat(input float64) float64 {
	return math.Round(input*10) / 10
}

func (t *TuyaTRV) forceMode() ForceMode {
	if (t.HeatingSetpoint - t.LocalTemperature) >= forceModeDiff {
		return ForceModeOpen
	}

	return ForceModeNormal
}

func (t *TuyaTRV) updateForceMode() error {
	f := t.forceMode()
	if t.Force == f {
		return nil
	}

	s := struct {
		Force ForceMode `json:"force"`
	}{Force: f}

	return t.write(s)
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

	s := struct {
		SetPoint float64 `json:"current_heating_setpoint"`
	}{SetPoint: formatFloat(temperature)}

	return t.write(s)
}

// TemperatureMode implements the TemperatureController interface
func (t *TuyaTRV) TemperatureMode() (components.TemperatureMode, error) {
	return components.TemperatureModeAuto, components.ErrNotImplemented
}

// SetTemperatureMode implements the TemperatureController interface
func (t *TuyaTRV) SetTemperatureMode(components.TemperatureMode) error {
	return components.ErrNotImplemented
}

// TemperatureCalibration implements the TemperatureController interface
func (t *TuyaTRV) TemperatureCalibration() (float64, error) {
	return t.LocalTemperatureCalibration, nil
}

// SetTemperatureCalibration implements the TemperatureController interface
func (t *TuyaTRV) SetTemperatureCalibration(c float64) error {
	s := struct {
		Calibration float64 `json:"local_temperature_calibration"`
	}{Calibration: formatFloat(c)}

	return t.write(s)
}
