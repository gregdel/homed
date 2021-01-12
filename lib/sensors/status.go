package sensors

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	register(TypeStatus, NewStatus)
}

// Status is a sensor that reports the status of a device
type Status struct {
	baseSensor

	Online bool `json:"online"`
}

// NewStatus returns a new status sensor
func NewStatus() Sensor {
	return &Status{}
}

// Type implements the Sensor interface
func (s *Status) Type() Type {
	return TypeStatus
}

// Collectors implements the Sensor interface
func (s *Status) Collectors(labels prometheus.Labels) []prometheus.Collector {
	return []prometheus.Collector{
		prometheus.NewGaugeFunc(
			prometheus.GaugeOpts{
				Name:        "homed_device_status",
				ConstLabels: labels,
			},
			func() float64 {
				if s.Online {
					return 1
				}
				return 0
			},
		),
	}
}

// Update implements the Sensor interface
func (s *Status) Update(value []byte) error {
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
