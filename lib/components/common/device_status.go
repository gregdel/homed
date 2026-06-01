package common

import (
	"bytes"
	"sync"
	"time"

	"github.com/gregdel/homed/lib/components"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	components.Register(components.TypeDeviceStatus, NewDeviceStatus)
}

// DeviceStatus is a component that reports the status of a device
type DeviceStatus struct {
	Component
	mu     sync.RWMutex
	online bool
}

type DeviceStatusSnapshot struct {
	ID           string             `json:"id"`
	UpdatedAt    *time.Time         `json:"updated_at"`
	FriendlyName string             `json:"friendly_name"`
	Hide         bool               `json:"hide"`
	Device       *components.Device `json:"device"`
	Online       bool               `json:"online"`
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
				if ds.IsOnline() {
					return 1
				}
				return 0
			},
		)}
}

// Update implements the Component interface
func (ds *DeviceStatus) Update(value []byte) error {
	online := bytes.Contains(bytes.ToLower(value), []byte("online"))
	ds.setOnline(online)

	if ds.Device() == nil {
		return components.ErrMissingDevice
	}

	// Update the device state
	ds.Device().SetOnline(online)

	return nil
}

func (ds *DeviceStatus) IsOnline() bool {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	return ds.online
}

func (ds *DeviceStatus) setOnline(online bool) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	ds.online = online
}

func (ds *DeviceStatus) Snapshot() any {
	return DeviceStatusSnapshot{
		ID:           ds.ID(),
		UpdatedAt:    ds.UpdatedAt.Load(),
		FriendlyName: ds.FriendlyName(),
		Hide:         ds.Hide,
		Device:       ds.Device(),
		Online:       ds.IsOnline(),
	}
}
