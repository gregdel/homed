package homed

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

// SensorTemperature is a sensor that handles temperatures
type SensorTemperature struct {
	BaseSensor
	Value float64
}

// NewSensorTemperature returns a new temperature sensor
func NewSensorTemperature() *SensorTemperature {
	return &SensorTemperature{}
}

// Type implements the Sensor interface
func (s *SensorTemperature) Type() SensorType {
	return SensorTypeTemperature
}

// Init implements the Sensor interface
func (s *SensorTemperature) Init() error {
	if s.device == nil {
		return ErrMissingDevice
	}

	c := prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "homed_temperature",
			ConstLabels: prometheus.Labels{
				"device": string(s.device.Name),
			},
		},
		func() float64 { return s.Value },
	)

	return prometheus.Register(c)
}

// Update implements the Sensor interface
func (s *SensorTemperature) Update(value []byte) error {
	v, err := strconv.ParseFloat(string(value), 64)
	if err != nil {
		return err
	}
	s.Value = v
	return nil
}
