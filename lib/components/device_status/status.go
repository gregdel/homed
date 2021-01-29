package status

import (
	"github.com/gregdel/homed/lib/components"
	"github.com/gregdel/homed/lib/components/common"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	components.Register(components.TypeDeviceStatus, NewDeviceStatus)
}

// DeviceStatus is a component that reports the status of a device
type DeviceStatus struct {
	common.Component
	common.DeviceStatus
}

// NewDeviceStatus returns a new status component
func NewDeviceStatus() components.Component {
	return &DeviceStatus{}
}

// Type implements the Component interface
func (s *DeviceStatus) Type() components.Type {
	return components.TypeDeviceStatus
}

// Collectors implements the Component interface
func (s *DeviceStatus) Collectors(labels prometheus.Labels) []prometheus.Collector {
	return []prometheus.Collector{
		prometheus.NewGaugeFunc(
			prometheus.GaugeOpts{
				Name:        "homed_device_status",
				ConstLabels: labels,
			},
			func() float64 {
				if s.IsAvailable() {
					return 1
				}
				return 0
			},
		),
	}
}
