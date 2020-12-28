package homed

import (
	"fmt"
	"strconv"

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
type Sensor struct {
	Device        *Device              `yaml:"-"`
	Name          string               `yaml:"name"`
	Type          SensorType           `yaml:"type"`
	Value         float64              `yaml:"-"`
	PromCollector prometheus.Collector `yaml:"-"`
}

// Update updates the sensor
func (s *Sensor) Update(value string) error {
	if s.PromCollector == nil {
		if s.Device == nil {
			return ErrMissingDevice
		}

		s.PromCollector = prometheus.NewGaugeFunc(prometheus.GaugeOpts{
			Name: "homed_" + string(s.Type),
			ConstLabels: prometheus.Labels{
				"device": string(s.Device.Name),
				"sensor": string(s.Name),
			},
		},
			func() float64 { return s.Value },
		)

		if err := prometheus.Register(s.PromCollector); err != nil {
			return err
		}
	}

	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return err
	}
	s.Value = v

	return nil
}

// Duplicate creates a copy of the sensor
func (s *Sensor) Duplicate() *Sensor {
	return &Sensor{
		Name: s.Name,
		Type: s.Type,
	}
}

func (s *Sensor) configFormat() interface{} {
	return s
}

func (s *Sensor) toConfig() (interface{}, error) {
	return s, nil
}

func (s *Sensor) fromConfig(data interface{}, homed *Homed) error {
	config, ok := data.(*Sensor)
	if !ok {
		return ErrInvalidConfigFormat
	}

	s = config
	return nil
}

// FileName implements the File interface
func (s *Sensor) FileName() string {
	return s.Name
}

// FileType implements the File interface
func (s *Sensor) FileType() FileType {
	return FileTypeSensor
}

// MQTTTopic returns the MQTT topic of the sensor
func (s *Sensor) MQTTTopic() (string, error) {
	if s.Device == nil {
		return "", ErrMissingDevice
	}

	return fmt.Sprintf(
		"home/devices/%s/sensor/%s/state",
		s.Device.Name,
		string(s.Type),
	), nil
}

// NewSensor returns a new sensor
func NewSensor(name string) *Sensor {
	return &Sensor{Name: name}
}
