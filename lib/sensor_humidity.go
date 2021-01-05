package homed

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

// SensorHumidity is a sensor that handles temperatures
type SensorHumidity struct {
	BaseSensor
	Value float64
}

// NewSensorHumidity returns a new humidity sensor
func NewSensorHumidity() *SensorHumidity {
	return &SensorHumidity{}
}

// Type implements the Sensor interface
func (s *SensorHumidity) Type() SensorType {
	return SensorTypeHumidity
}

// Init implements the Sensor interface
func (s *SensorHumidity) Init() error {
	if s.device == nil {
		return ErrMissingDevice
	}

	c := prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "homed_humidity",
			ConstLabels: prometheus.Labels{
				"device": string(s.device.Name),
			},
		},
		func() float64 { return s.Value },
	)

	return prometheus.Register(c)
}

// Update implements the Sensor interface
func (s *SensorHumidity) Update(value []byte) error {
	v, err := strconv.ParseFloat(string(value), 64)
	if err != nil {
		return err
	}
	s.Value = v
	return nil
}
