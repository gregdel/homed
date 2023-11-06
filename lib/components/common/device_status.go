package common

import (
	"fmt"
	"strings"

	"github.com/gregdel/homed/lib/components"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/atomic"
)

func init() {
	components.Register(components.TypeDeviceStatus, NewDeviceStatus)
}

// DeviceStatus is a component that reports the status of a device
type DeviceStatus struct {
	Component
	Online atomic.Bool `json:"online"`
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
				if ds.Online.Load() {
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
		ds.Online.Store(true)
	case "offline":
		ds.Online.Store(false)
	default:
		return fmt.Errorf("device status: invalid component status: %s", v)
	}

	if ds.Device() == nil {
		return components.ErrMissingDevice
	}

	// Update the device state
	ds.Device().Online.Store(ds.Online.Load())

	return nil
}
