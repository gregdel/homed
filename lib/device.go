package homed

import (
	"encoding/json"
	"fmt"

	"github.com/gregdel/homed/lib/sensors"
	"github.com/prometheus/client_golang/prometheus"
)

// Device represents a device
type Device struct {
	Room    *Room
	Name    string
	Sensors []sensors.Sensor
	Actions []Action
}

// NewDevice creates a new device
func NewDevice(name string) *Device {
	return &Device{Name: name}
}

// MarshalJSON implements the json.Marshaler interface
func (d *Device) MarshalJSON() ([]byte, error) {
	type sensorWithType struct {
		sensors.Sensor `json:"values"`
		Type           string `json:"type"`
	}

	s := make([]sensorWithType, len(d.Sensors))
	for i := 0; i < len(d.Sensors); i++ {
		s[i] = sensorWithType{
			Sensor: d.Sensors[i],
			Type:   string(d.Sensors[i].Type()),
		}
	}

	out := struct {
		Name    string           `json:"name"`
		Sensors []sensorWithType `json:"sensors"`
	}{
		Name:    d.Name,
		Sensors: s,
	}

	return json.Marshal(out)
}

// AddSensor adds a sensor to the device
func (d *Device) AddSensor(sensorType, topic string) (sensors.Sensor, error) {
	sensor, err := sensors.New(sensorType)
	if err != nil {
		return nil, err
	}

	if d.Sensors == nil {
		d.Sensors = []sensors.Sensor{}
	}

	labels := prometheus.Labels{"device": d.Name}
	if d.Room != nil {
		labels["room"] = string(d.Room.Name)
	}

	collectors := sensor.Collectors(labels)
	for _, c := range collectors {
		if err := prometheus.Register(c); err != nil {
			return nil, err
		}

	}

	d.Sensors = append(d.Sensors, sensor)

	return sensor, nil
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
