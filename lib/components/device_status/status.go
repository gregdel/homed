package status

import (
	"fmt"
	"strings"

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
	Online bool `json:"online"`
}

// NewDeviceStatus returns a new status component
func NewDeviceStatus() components.Component {
	return &DeviceStatus{}
}

// Type implements the Component interface
func (ds *DeviceStatus) Type() components.Type {
	return components.TypeDeviceStatus
}

// Collectors implements the Component interface
func (ds *DeviceStatus) Collectors(labels prometheus.Labels) []prometheus.Collector {
	return []prometheus.Collector{
		components.GaugeCollector("device_status", labels,
			func() float64 {
				if ds.Online {
					return 1
				}
				return 0
			},
		)}
}

// Update implements the Component interface
func (ds *DeviceStatus) Update(value []byte) error {
	v := strings.ToLower(string(value))
	switch v {
	case "online":
		ds.Online = true
	case "offline":
		ds.Online = false
	default:
		return fmt.Errorf("device status: invalid component status: %s", v)
	}

	if ds.Device() == nil {
		return components.ErrMissingDevice
	}

	// Update the device state
	ds.Device().Online = ds.Online

	return nil
}
