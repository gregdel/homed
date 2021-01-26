package status

import (
	"fmt"
	"strings"

	"github.com/gregdel/homed/lib/components"
	base "github.com/gregdel/homed/lib/components/base_component"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	components.Register(components.TypeDeviceStatus, NewDeviceStatus)
}

// DeviceStatus is a component that reports the status of a device
type DeviceStatus struct {
	base.Component

	Online bool `json:"online"`
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
				if s.Online {
					return 1
				}
				return 0
			},
		),
	}
}

// Update implements the Component interface
func (s *DeviceStatus) Update(value []byte) error {
	v := strings.ToLower(string(value))
	switch v {
	case "online":
		s.Online = true
	case "offline":
		s.Online = false
	default:
		return fmt.Errorf("components: invalid component status: %s", v)
	}
	return nil
}
