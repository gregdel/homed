package components

import (
	"github.com/gregdel/homed/lib/schedule"
)

// Switch is an interface to reprensents a switch
type Switch interface {
	IsOn() bool
	Toggle() error
	SetOn() error
	SetOff() error
	Set(state bool) error
}

// Publisher is an interface to publish the component state
type Publisher interface {
	PublishState() error
}

// Scheduled is an interface to handle the schedule of a component
type Scheduled interface {
	Schedule() *schedule.Schedule
	SetSchedule(*schedule.Schedule)
	SaveSchedule(path string) error
	LoadSchedule(path string) error
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
