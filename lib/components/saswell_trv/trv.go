package trv

import (
	"encoding/json"
	"math"

	"github.com/gregdel/homed/lib/components"
	"github.com/gregdel/homed/lib/components/common"
	"github.com/prometheus/client_golang/prometheus"
)

// Make sure that the module is a temperature controller
var _ components.TemperatureController = (*SaswellTRV)(nil)

func init() {
	components.Register(components.TypeSaswellTRV, New)
}

// SystemMode represents the system modes
type SystemMode string

// Force modes
var (
	SystemModeAuto SystemMode = "auto"
	SystemModeHeat SystemMode = "heat"
	SystemModeOff  SystemMode = "off"
)

// SaswellTRV is a component that handles temperatures
type SaswellTRV struct {
	common.Component

	HeatingSetpoint             float64    `json:"current_heating_setpoint"`
	LocalTemperatureCalibration float64    `json:"local_temperature_calibration"`
	LocalTemperature            float64    `json:"local_temperature"`
	BatteryLow                  bool       `json:"battery_low"`
	Mode                        SystemMode `json:"system_mode"`
}

// New returns a new component for Saswell TRVs
func New() components.Component {
	return &SaswellTRV{}
}

// Type implements the Component interface
func (t *SaswellTRV) Type() components.Type {
	return components.TypeSaswellTRV
}

// Collectors implements the Component interface
func (t *SaswellTRV) Collectors(labels prometheus.Labels) []prometheus.Collector {
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

func (t *SaswellTRV) write(input interface{}) error {
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

// Update implements the Component interface
func (t *SaswellTRV) Update(value []byte) error {
	return json.Unmarshal(value, t)
}

// SetTemperature implements the TemperatureGetterSetter interface
func (t *SaswellTRV) SetTemperature(float64) error {
	return components.ErrNotImplemented
}

// Temperature implements the TemperatureGetter interface
func (t *SaswellTRV) Temperature() (float64, error) {
	return t.LocalTemperature, nil
}

// TemperatureTarget implements the TemperatureController interface
func (t *SaswellTRV) TemperatureTarget() (float64, error) {
	return t.HeatingSetpoint, nil
}

// SetTemperatureTarget implements the TemperatureController interface
func (t *SaswellTRV) SetTemperatureTarget(temperature float64) error {
	if temperature == t.HeatingSetpoint {
		return nil
	}

	s := struct {
		SetPoint float64 `json:"current_heating_setpoint"`
	}{SetPoint: formatFloat(temperature)}

	return t.write(s)
}

// TemperatureMode implements the TemperatureController interface
func (t *SaswellTRV) TemperatureMode() (components.TemperatureMode, error) {
	return components.TemperatureModeAuto, components.ErrNotImplemented
}

// SetTemperatureMode implements the TemperatureController interface
func (t *SaswellTRV) SetTemperatureMode(components.TemperatureMode) error {
	return components.ErrNotImplemented
}

// TemperatureCalibration implements the TemperatureController interface
func (t *SaswellTRV) TemperatureCalibration() (float64, error) {
	return t.LocalTemperatureCalibration, nil
}

// SetTemperatureCalibration implements the TemperatureController interface
func (t *SaswellTRV) SetTemperatureCalibration(c float64) error {
	s := struct {
		Calibration float64 `json:"local_temperature_calibration"`
	}{Calibration: formatFloat(c)}

	return t.write(s)
}
