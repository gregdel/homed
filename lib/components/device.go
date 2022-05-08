package components

import (
	"fmt"
	"sync"
)

// Device represents a device
type Device struct {
	mu sync.Mutex

	Name       string               `json:"name"`
	Room       string               `json:"room"`
	Online     bool                 `json:"online"`
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

// String implements the Stringer interface
func (d *Device) String() string {
	var out string
	out += fmt.Sprintf("Device: %s\n", d.Name)
	out += fmt.Sprintf("Room:   %s\n", d.Room)
	out += fmt.Sprintf("Online: %t\n", d.Online)
	out += fmt.Sprintf("Components:\n")
	for n := range d.Components {
		out += fmt.Sprintf("  - %s\n", n)
	}
	return out
}
