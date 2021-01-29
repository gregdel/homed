package common

import (
	"fmt"
	"strings"
)

// DeviceStatus represents a generic sensor
type DeviceStatus struct {
	Component

	Available bool `json:"on"`
}

// IsAvailable implements the Status interface
func (ds *DeviceStatus) IsAvailable() bool {
	return ds.Available
}

// Update implements the Component interface
func (ds *DeviceStatus) Update(value []byte) error {
	v := strings.ToLower(string(value))
	switch v {
	case "online":
		ds.Available = true
	case "offline":
		ds.Available = false
		return fmt.Errorf("device status: invalid component status: %s", v)
	}
	return nil
}
