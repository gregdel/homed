package sensors

import (
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	register(TypeHumidity, NewHumidity)
}

// Humidity is a sensor that handles temperatures
type Humidity struct {
	baseSensor
	float64Sensor
}

// NewHumidity returns a new humidity sensor
func NewHumidity() Sensor {
	return &Humidity{}
}

// Type implements the Sensor interface
func (s *Humidity) Type() Type {
	return TypeHumidity
}

// Collectors implements the Sensor interface
func (s *Humidity) Collectors(labels prometheus.Labels) []prometheus.Collector {
	return singleCollector(s, labels, func() float64 { return s.Value })
}
