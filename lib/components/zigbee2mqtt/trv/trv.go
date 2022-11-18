package trv

import (
	"encoding/json"
	"math"
	"time"

	"github.com/gregdel/homed/lib/components"
	"github.com/gregdel/homed/lib/components/common"
	"github.com/prometheus/client_golang/prometheus"
)

// Make sure that the module is a temperature controller
var _ components.TemperatureController = (*TRV)(nil)

func init() {
	components.Register(components.TypeZigbeeTRV, New)
}

var (
	sasswellOffset            = 2.0
	forceModeDiff             = 1.0
	calibrationRequestTimeout = 30 * time.Minute
)

// TRV represents a zigbee2mqtt TRV
type TRV struct {
	common.Component

	// Track when the calibration was requested
	CalibrationRequestTime *time.Time `json:"calibration_request_time"`

	HeatingSetpoint             float64    `json:"current_heating_setpoint"`
	LocalTemperatureCalibration float64    `json:"local_temperature_calibration"`
	LocalTemperature            float64    `json:"local_temperature"`
	BatteryLow                  bool       `json:"battery_low"`
	Mode                        SystemMode `json:"system_mode"`
	Force                       ForceMode  `json:"force"`
	Position                    float64    `json:"position"`
}

// New returns a new component for Tuya TRVs
func New() components.Component {
	return &TRV{
		Position: -1,
	}
}

// Type implements the Component interface
func (t *TRV) Type() components.Type {
	return components.TypeZigbeeTRV
}

// Collectors implements the Component interface
func (t *TRV) Collectors(labels prometheus.Labels) []prometheus.Collector {
	return []prometheus.Collector{
		components.GaugeCollector("temperature", labels,
			func() float64 { return t.LocalTemperature },
		),
		components.GaugeCollector("temperature_calibration", labels,
			func() float64 { return t.LocalTemperatureCalibration },
		),
		components.GaugeCollector("heating_set_point", labels,
			func() float64 { return t.HeatingSetpoint },
		),
		components.GaugeCollector("trv_position", labels,
			func() float64 { return t.Position },
		),
	}
}

// This is where the ugly shit begins
func (t *TRV) isSaswell() bool {
	return t.Force == ForceModeUnavailable
}

// Update implements the Component interface
func (t *TRV) Update(value []byte) error {
	oldTemperature := t.LocalTemperature

	if err := json.Unmarshal(value, t); err != nil {
		return err
	}

	// The temperature was updated, let's assume it take the calibration into
	// account
	if t.LocalTemperature != oldTemperature {
		t.CalibrationRequestTime = nil
	}

	// After some time, the calibration might not be relevant anymore
	if (t.CalibrationRequestTime != nil) &&
		time.Since(*t.CalibrationRequestTime) > calibrationRequestTimeout {
		t.CalibrationRequestTime = nil
	}

	return nil
}

// PostUpdate implements the Component interface
func (t *TRV) PostUpdate() error {
	err := t.Component.PostUpdate()
	if err != nil {
		return err
	}

	return t.updateMode()
}

// Temperature implements the TemperatureGetter interface
func (t *TRV) Temperature() (float64, error) {
	if !t.Device().Online {
		return 0, components.ErrDeviceOffline
	}

	return t.LocalTemperature, nil
}

// TemperatureTarget implements the TemperatureController interface
func (t *TRV) TemperatureTarget() (float64, error) {
	if t.isSaswell() {
		return t.HeatingSetpoint - sasswellOffset, nil
	}

	return t.HeatingSetpoint, nil
}

// SetTemperatureTarget implements the TemperatureController interface
func (t *TRV) SetTemperatureTarget(temperature float64) error {
	sp, _ := t.TemperatureTarget()
	if temperature == sp {
		return nil
	}

	if t.isSaswell() {
		temperature += sasswellOffset
	}

	s := struct {
		SetPoint float64 `json:"current_heating_setpoint"`
	}{SetPoint: formatFloat(temperature)}

	return t.write(s)
}

// TemperatureCalibration implements the TemperatureController interface
func (t *TRV) TemperatureCalibration() (float64, error) {
	if !t.Device().Online {
		return 0, components.ErrDeviceOffline
	}

	if t.CalibrationRequestTime != nil {
		return 0, components.ErrOperatingInProgress
	}

	return t.LocalTemperatureCalibration, nil
}

// SetTemperatureCalibration implements the TemperatureController interface
func (t *TRV) SetTemperatureCalibration(c float64) error {
	if t.CalibrationRequestTime != nil {
		return components.ErrOperatingInProgress
	}

	// Saswell TRV doesn't support floats, let's only use ints for now
	s := struct {
		Calibration float64 `json:"local_temperature_calibration"`
	}{Calibration: math.Round(c)}

	if err := t.write(s); err != nil {
		return err
	}

	now := time.Now()
	t.CalibrationRequestTime = &now
	return nil
}

// Format the float to 1 digit of precision
func formatFloat(input float64) float64 {
	return math.Round(input*10) / 10
}

func (t *TRV) write(data interface{}) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return t.WriteCommand(payload)
}

func (t *TRV) updateMode() error {
	if t.Force == ForceModeUnavailable {
		return nil
	}

	newMode := ForceModeNormal
	if (t.HeatingSetpoint - t.LocalTemperature) >= forceModeDiff {
		newMode = ForceModeOpen
	}

	if t.Force == newMode {
		return nil
	}

	s := struct {
		Force ForceMode `json:"force"`
	}{Force: newMode}

	return t.write(s)
}
