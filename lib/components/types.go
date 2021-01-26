package components

// TemperatureGetter is an interface to reprensents something that holds a
// temperature
type TemperatureGetter interface {
	Component
	GetTemperature() (float64, error)
}

// TemperatureSetter is an interface to reprensents something that sets a
// temperature
type TemperatureSetter interface {
	Component
	SetTemperature(float64) error
}
