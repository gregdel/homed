package homed

import (
	"github.com/gregdel/homed/lib/components"
	"github.com/prometheus/client_golang/prometheus"
)

// Device represents a device
type Device struct {
	Room       *Room                 `json:"-"`
	Name       string                `json:"name"`
	Components components.Components `json:"components"`
}

// NewDevice creates a new device
func NewDevice(name string) *Device {
	return &Device{
		Name:       name,
		Components: components.New(),
	}
}

// AddComponent adds a component to the device
func (d *Device) AddComponent(cfg components.Config) (components.Component, error) {
	labels := prometheus.Labels{"device": d.Name}
	if d.Room != nil {
		labels["room"] = string(d.Room.Name)
	}

	return d.Components.Add(cfg, labels)
}
