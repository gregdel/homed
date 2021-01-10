package sensors

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	register(TypeTemperature, NewTemperature)
}

// Temperature is a sensor that handles temperatures
type Temperature struct {
	BaseSensor
	Value float64
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
	return []prometheus.Collector{
		prometheus.NewGaugeFunc(
			prometheus.GaugeOpts{
				Name:        "homed_temperature",
				ConstLabels: labels,
			},
			func() float64 { return s.Value },
		),
	}
}

// Update implements the Sensor interface
func (s *Temperature) Update(value []byte) error {
	v, err := strconv.ParseFloat(string(value), 64)
	if err != nil {
		return err
	}
	s.Value = v
	return nil
}
