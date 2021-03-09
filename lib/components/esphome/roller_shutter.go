package esphome

import (
	"github.com/gregdel/homed/lib/components"
	"github.com/gregdel/homed/lib/components/common"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	components.Register(components.TypeRollerShutter, NewRollerShutter)
}

// RollerShutter is a component that controls a cover
type RollerShutter struct {
	common.RollerShutter
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
