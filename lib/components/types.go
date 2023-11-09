package components

// Switch is an interface to reprensents a switch
type Switch interface {
	IsOn() bool
	Toggle() error
	TurnOn() error
	TurnOff() error
}

// Sensor reprensents a sensor that holds a value
type Sensor interface {
	SensorValue() float64
}

// BinarySensor reprensents a binary sensor that holds a true/false value
type BinarySensor interface {
	IsOn() bool
}

// RollerShutter is an interface to reprensents a roller shutter
type RollerShutter interface {
	Open() error
	Close() error
	Stop() error
	IsOpen() bool
	IsClosed() bool
	OpenedAt() float64
}
