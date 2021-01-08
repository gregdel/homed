package homed

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

// SensorStatus is a sensor that reports the status of a device
type SensorStatus struct {
	BaseSensor
	Online bool
}

// NewSensorStatus returns a new status sensor
func NewSensorStatus() *SensorStatus {
	return &SensorStatus{}
}

// Type implements the Sensor interface
func (s *SensorStatus) Type() SensorType {
	return SensorTypeStatus
}

// Init implements the Sensor interface
func (s *SensorStatus) Init() error {
	if s.device == nil {
		return ErrMissingDevice
	}

	c := prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "homed_device_status",
			ConstLabels: prometheus.Labels{
				"device": string(s.device.Name),
			},
		},
		func() float64 {
			if s.Online {
				return 1
			}
			return 0
		},
	)

	return prometheus.Register(c)
}

// Update implements the Sensor interface
func (s *SensorStatus) Update(value []byte) error {
	v := string(value)
	switch v {
	case "online":
		s.Online = true
	case "offline":
		s.Online = false
	default:
		return fmt.Errorf("homed: invalid sensor status: %s", v)
	}
	return nil
}
