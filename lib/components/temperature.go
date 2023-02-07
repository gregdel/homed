package components

import "time"

// TemperatureMode represents a temperature control mode
type TemperatureMode string

// Available modes
var (
	TemperatureModeAuto            TemperatureMode = "auto"
	TemperatureModeFixed           TemperatureMode = "fixed"
	TemperatureModeDuration        TemperatureMode = "duration"
	TemperatureModeUntilDate       TemperatureMode = "until_date"
	TemperatureModeUntilNextChange TemperatureMode = "until_next_change"
)

// TemperatureGetter is an interface to reprensents something that holds a
// temperature
type TemperatureGetter interface {
	Component
	Temperature() (float64, error)
}

// TemperatureController is an interface to reprensents something that be get
// or set
type TemperatureController interface {
	TemperatureGetter

	TemperatureTarget() (float64, error)
	SetTemperatureTarget(float64) error

	TemperatureCalibration() (float64, error)
	SetTemperatureCalibration(float64) error
}

// TemperatureControllerInternal is an interface to reprensents something that be get
// or set
type TemperatureControllerInternal interface {
	Publisher
	Scheduled
	TemperatureController

	SetTemperature(float64) error

	TemperatureMode() (TemperatureMode, error)
	SetTemperatureMode(TemperatureMode) error

	SetTemperatureManualTarget(float64) error
	TemperatureManualTarget() (float64, error)

	TemperatureModeManualUntil() (*time.Time, error)
	SetTemperatureModeManualUntil(*time.Time) error
}
