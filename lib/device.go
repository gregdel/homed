package homed

import "fmt"

// Device represents a device
type Device struct {
	Room    *Room
	Name    string
	Sensors []Sensor
}

// NewDevice creates a new device
func NewDevice(name string) *Device {
	return &Device{Name: name}
}

// AddSensor adds a sensor to the device
func (d *Device) AddSensor(sensorType, topic string) (Sensor, error) {
	var sensor Sensor
	switch sensorType {
	case "temperature":
		sensor = NewSensorTemperature()
	case "humidity":
		sensor = NewSensorHumidity()
	case "wifi_signal":
		sensor = NewSensorWifiSignal()
	default:
		return nil, fmt.Errorf("homed: invalid sensor type: %s", sensorType)
	}
	sensor.SetMQTTTopic(topic)

	if d.Sensors == nil {
		d.Sensors = []Sensor{}
	}
	d.Sensors = append(d.Sensors, sensor)

	return sensor, nil
}
