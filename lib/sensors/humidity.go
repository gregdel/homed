package sensors

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	register(TypeHumidity, NewHumidity)
}

// Humidity is a sensor that handles temperatures
type Humidity struct {
	BaseSensor
	Value float64
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
	return []prometheus.Collector{
		prometheus.NewGaugeFunc(
			prometheus.GaugeOpts{
				Name:        "homed_humidity",
				ConstLabels: labels,
			},
			func() float64 { return s.Value },
		),
	}
}

// Update implements the Sensor interface
func (s *Humidity) Update(value []byte) error {
	v, err := strconv.ParseFloat(string(value), 64)
	if err != nil {
		return err
	}
	s.Value = v
	return nil
}
