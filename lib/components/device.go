package components

import (
	"fmt"
	"sync"

	"go.uber.org/atomic"
)

// Device represents a device
type Device struct {
	mu sync.Mutex

	Name       string               `json:"name"`
	Room       string               `json:"room"`
	Online     atomic.Bool          `json:"online"`
	Components map[string]Component `json:"-"`
}

// NewDevice returns a new device
func NewDevice(name, room string) *Device {
	return &Device{
		Name:       name,
		Room:       room,
		Components: map[string]Component{},
	}
}

// AddComponent adds a component on the device
func (d *Device) AddComponent(c Component) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, ok := d.Components[c.ID()]; ok {
		return fmt.Errorf("devices: component %s already added", c.ID())
	}

	d.Components[c.ID()] = c
	return nil
}

// IsOnline tells if the device is online or not
func (d *Device) IsOnline() bool {
	return d.Online.Load()
}
