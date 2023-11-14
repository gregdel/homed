package components

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

// TemperatureGetter is the interface implemented by anything that can return a
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
	Component
	Scheduled

	IsHeating() bool
}

// HumidityGetter is the interface implemented by anything that can return a
// humidity
type HumidityGetter interface {
	Component
	Humidity() (float64, error)
}
