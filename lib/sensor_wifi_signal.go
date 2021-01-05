package homed

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

// SensorWifiSignal is a sensor that handles temperatures
type SensorWifiSignal struct {
	BaseSensor
	Value float64
}

// NewSensorWifiSignal returns a new wifi signal sensor
func NewSensorWifiSignal() *SensorWifiSignal {
	return &SensorWifiSignal{}
}

// Type implements the Sensor interface
func (s *SensorWifiSignal) Type() SensorType {
	return SensorTypeWifiSignal
}

// Init implements the Sensor interface
func (s *SensorWifiSignal) Init() error {
	if s.device == nil {
		return ErrMissingDevice
	}

	c := prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "homed_wifi_signal",
			ConstLabels: prometheus.Labels{
				"device": string(s.device.Name),
			},
		},
		func() float64 { return s.Value },
	)

	return prometheus.Register(c)
}

// Update implements the Sensor interface
func (s *SensorWifiSignal) Update(value []byte) error {
	v, err := strconv.ParseFloat(string(value), 64)
	if err != nil {
		return err
	}
	s.Value = v
	return nil
}
