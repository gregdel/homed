package homed

import (
	"github.com/prometheus/client_golang/prometheus"
)

// SensorType reprensents a sensor type
type SensorType string

// SensorTypes
var (
	SensorTypeUnknown     SensorType = "unknown"
	SensorTypeTemperature SensorType = "temperature"
	SensorTypeHumidity    SensorType = "humidity"
	SensorTypeWifiSignal  SensorType = "wifi_signal"
)

// Sensor represents a sensor
type Sensor interface {
	MQTTTopic() string
	SetDevice(*Device)
	SetMQTTTopic(string)
	Type() SensorType
	Update(string) error
}

// BaseSensor represents a basic sensor
type BaseSensor struct {
	device        *Device
	mqttTopic     string
	promCollector prometheus.Collector `yaml:"-"`
}

// MQTTTopic implements the Sensor interface
func (s *BaseSensor) MQTTTopic() string {
	return s.mqttTopic
}

// SetMQTTTopic implements the Sensor interface
func (s *BaseSensor) SetMQTTTopic(value string) {
	s.mqttTopic = value
}

// Type implements the Sensor interface
func (s *BaseSensor) Type() SensorType {
	return SensorTypeUnknown
}

// Update implements the Sensor interface
func (s *BaseSensor) Update(_ string) error {
	return nil
}

// SetDevice implements the Sensor interface
func (s *BaseSensor) SetDevice(d *Device) {
	s.device = d
}
