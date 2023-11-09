package trv

import (
	"encoding/json"
	"math"
	"sync"
	"time"

	"github.com/gregdel/homed/lib/components"
	"github.com/gregdel/homed/lib/components/common"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/atomic"
)

var zeroTime time.Time

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

// Data represents the data received from the TRV
type Data struct {
	HeatingSetpoint             atomic.Float64 `json:"current_heating_setpoint"`
	LocalTemperatureCalibration atomic.Float64 `json:"local_temperature_calibration"`
	LocalTemperature            atomic.Float64 `json:"local_temperature"`
	BatteryLow                  atomic.Bool    `json:"battery_low"`
	Mode                        atomic.String  `json:"system_mode"`
	Force                       atomic.String  `json:"force"`
	Position                    atomic.Float64 `json:"position"`
}

// TRV represents a zigbee2mqtt TRV
type TRV struct {
	mu sync.RWMutex

	common.Component
	Data

	// Track when the calibration was requested
	CalibrationRequestTime atomic.Time `json:"calibration_request_time"`
}

// New returns a new component for Tuya TRVs
func New() components.Component {
	trv := &TRV{}
	trv.Position.Store(-1)
	return trv
}

// Type implements the Component interface
func (t *TRV) Type() components.Type {
	return components.TypeZigbeeTRV
}

// Collectors implements the Component interface
func (t *TRV) Collectors(labels prometheus.Labels) []prometheus.Collector {
	return []prometheus.Collector{
		components.GaugeCollector("temperature", labels,
			func() float64 { return t.LocalTemperature.Load() },
		),
		components.GaugeCollector("temperature_calibration", labels,
			func() float64 { return t.LocalTemperatureCalibration.Load() },
		),
		components.GaugeCollector("heating_set_point", labels,
			func() float64 { return t.HeatingSetpoint.Load() },
		),
		components.GaugeCollector("trv_position", labels,
			func() float64 { return t.Position.Load() },
		),
	}
}

// This is where the ugly shit begins
func (t *TRV) isSaswell() bool {
	return t.Force.Load() == ForceModeUnavailable
}

// Update implements the Component interface
func (t *TRV) Update(value []byte) error {
	oldTemperature := t.LocalTemperature.Load()

	if err := json.Unmarshal(value, &t.Data); err != nil {
		return err
	}

	// The temperature was updated, let's assume it take the calibration into
	// account
	if t.LocalTemperature.Load() != oldTemperature {
		t.CalibrationRequestTime.Store(zeroTime)
	}

	// After some time, the calibration might not be relevant anymore
	requestTime := t.CalibrationRequestTime.Load()
	if !requestTime.IsZero() &&
		time.Since(requestTime) > calibrationRequestTimeout {
		t.CalibrationRequestTime.Store(zeroTime)
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
	if !t.Device().IsOnline() {
		return 0, components.ErrDeviceOffline
	}

	return t.LocalTemperature.Load(), nil
}

// TemperatureTarget implements the TemperatureController interface
func (t *TRV) TemperatureTarget() (float64, error) {
	heatingSetpoint := t.HeatingSetpoint.Load()
	if t.isSaswell() {
		return heatingSetpoint - sasswellOffset, nil
	}

	return heatingSetpoint, nil
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
	if !t.Device().IsOnline() {
		return 0, components.ErrDeviceOffline
	}

	if !t.CalibrationRequestTime.Load().IsZero() {
		return 0, components.ErrOperatingInProgress
	}

	return t.LocalTemperatureCalibration.Load(), nil
}

// SetTemperatureCalibration implements the TemperatureController interface
func (t *TRV) SetTemperatureCalibration(c float64) error {
	if !t.CalibrationRequestTime.Load().IsZero() {
		return components.ErrOperatingInProgress
	}

	// Saswell TRV doesn't support floats, let's only use ints for now
	s := struct {
		Calibration float64 `json:"local_temperature_calibration"`
	}{Calibration: math.Round(c)}

	if err := t.write(s); err != nil {
		return err
	}

	t.CalibrationRequestTime.Store(time.Now())
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

func (t *TRV) saswellEnsureHeat() error {
	if t.Mode.Load() == SystemModeHeat {
		return nil
	}

	s := struct {
		Mode string `json:"system_mode"`
	}{Mode: SystemModeHeat}

	return t.write(s)
}

func (t *TRV) tuyaForce() error {
	diff := (t.HeatingSetpoint.Load() - t.LocalTemperature.Load())
	newMode := ForceModeNormal
	if diff >= forceModeDiff {
		newMode = ForceModeOpen
	}

	if t.Mode.Load() == newMode {
		return nil
	}

	s := struct {
		Force string `json:"force"`
	}{Force: newMode}

	return t.write(s)
}

func (t *TRV) updateMode() error {
	if t.isSaswell() {
		// Ensure that the TRV is in "heat" mode and not auto
		return t.saswellEnsureHeat()
	}

	// Force the TRV to open
	return t.tuyaForce()
}
