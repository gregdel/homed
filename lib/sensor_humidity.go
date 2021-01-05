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

// Update implements the Sensor interface
func (s *SensorHumidity) Update(value string) error {
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
