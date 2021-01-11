package homed

import (
	"fmt"

	"github.com/gregdel/homed/lib/sensors"
	"github.com/prometheus/client_golang/prometheus"
)

// Device represents a device
type Device struct {
	Room    *Room           `json:"-"`
	Name    string          `json:"name"`
	Sensors sensors.Sensors `json:"sensors"`
	Actions []Action        `json:"-"`
}

// NewDevice creates a new device
func NewDevice(name string) *Device {
	return &Device{
		Name:    name,
		Sensors: sensors.New(),
	}
}

// AddSensor adds a sensor to the device
func (d *Device) AddSensor(sensorType string) (sensors.Sensor, error) {
	labels := prometheus.Labels{"device": d.Name}
	if d.Room != nil {
		labels["room"] = string(d.Room.Name)
	}

	return d.Sensors.Add(sensorType, labels)
}

// AddAction adds an action to the device
func (d *Device) AddAction(actionType, topic string) (Action, error) {
	var action Action
	switch actionType {
	case "set_temperature":
		action = NewActionSetTemperature()
	default:
		return nil, fmt.Errorf("homed: invalid action type: %s", actionType)
	}

	if d.Actions == nil {
		d.Actions = []Action{}
	}

	d.Actions = append(d.Actions, action)

	return action, nil
}
