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

// Update implements the Sensor interface
func (s *SensorTemperature) Update(value string) error {
	if s.promCollector == nil {
		if s.device == nil {
			return ErrMissingDevice
		}

		s.promCollector = prometheus.NewGaugeFunc(prometheus.GaugeOpts{
			Name: "homed_" + string(s.Type()),
			ConstLabels: prometheus.Labels{
				"device": string(s.device.Name),
			},
		},
			func() float64 { return s.Value },
		)

		if err := prometheus.Register(s.promCollector); err != nil {
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
