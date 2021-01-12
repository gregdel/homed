package sensors

import (
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	register(TypeTemperature, NewTemperature)
}

// Temperature is a sensor that handles temperatures
type Temperature struct {
	baseSensor
	float64Sensor
}

// NewTemperature returns a new temperature sensor
func NewTemperature() Sensor {
	return &Temperature{}
}

// Type implements the Sensor interface
func (s *Temperature) Type() Type {
	return TypeTemperature
}

// Collectors implements the Sensor interface
func (s *Temperature) Collectors(labels prometheus.Labels) []prometheus.Collector {
	return singleCollector(s, labels, func() float64 { return s.Value })
}
