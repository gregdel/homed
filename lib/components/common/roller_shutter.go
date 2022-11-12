package common

import (
	"github.com/gregdel/homed/lib/components"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	components.Register(components.TypeRollerShutter, NewRollerShutter)
}

// RollerShutter represents a generic roller shutter
type RollerShutter struct {
	GenericSensor
}

// NewRollerShutter returns a new cover component
func NewRollerShutter() components.Component {
	return &RollerShutter{}
}

// Type implements the Component interface
func (rs *RollerShutter) Type() components.Type {
	return components.TypeRollerShutter
}

// Collectors implements the Component interface
func (rs *RollerShutter) Collectors(labels prometheus.Labels) []prometheus.Collector {
	return nil
}

// OpenedAt implements the RollerShutter interface
func (rs *RollerShutter) OpenedAt() float64 {
	return rs.Value
}

// IsOpen implements the RollerShutter interface
func (rs *RollerShutter) IsOpen() bool {
	return rs.Value == 100
}

// IsClosed implements the RollerShutter interface
func (rs *RollerShutter) IsClosed() bool {
	return !rs.IsOpen()
}

// Open implements the RollerShutter interface
func (rs *RollerShutter) Open() error {
	return rs.WriteCommand([]byte("open"))
}

// Close implements the RollerShutter interface
func (rs *RollerShutter) Close() error {
	return rs.WriteCommand([]byte("close"))
}

// Stop implements the RollerShutter interface
func (rs *RollerShutter) Stop() error {
	return rs.WriteCommand([]byte("stop"))
}
