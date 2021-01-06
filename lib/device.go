package homed

import "fmt"

// Device represents a device
type Device struct {
	Room    *Room
	Name    string
	Sensors []Sensor
	Actions []Action
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
	case "zigbee2mqtt_tuya":
		sensor = NewSensorZigbee2MQTTTuya()
	default:
		return nil, fmt.Errorf("homed: invalid sensor type: %s", sensorType)
	}
	sensor.SetMQTTTopic(topic)
	sensor.SetDevice(d)

	if d.Sensors == nil {
		d.Sensors = []Sensor{}
	}

	if err := sensor.Init(); err != nil {
		return nil, err
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
