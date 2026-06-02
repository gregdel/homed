package components

import (
	"encoding/json"
	"fmt"
	"sync"
)

// Device represents a device
type Device struct {
	mu sync.Mutex

	Name                string `json:"name"`
	Room                string `json:"room"`
	online              bool
	availabilityTracked bool
	Components          map[string]Component `json:"-"`
}

type deviceJSON struct {
	Name   string `json:"name"`
	Room   string `json:"room"`
	Online bool   `json:"online"`
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
	d.mu.Lock()
	defer d.mu.Unlock()

	return d.online
}

func (d *Device) SetOnline(online bool) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.online = online
}

// SetAvailabilityTracked marks the device as having an explicit availability
// source.
func (d *Device) SetAvailabilityTracked() {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.availabilityTracked = true
}

// MetricsAvailable tells if telemetry metrics should be exposed.
func (d *Device) MetricsAvailable() bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	return !d.availabilityTracked || d.online
}

func (d *Device) MarshalJSON() ([]byte, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	return json.Marshal(deviceJSON{
		Name:   d.Name,
		Room:   d.Room,
		Online: d.online,
	})
}
