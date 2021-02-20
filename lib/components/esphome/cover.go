package esphome

import (
	"github.com/gregdel/homed/lib/components"
	"github.com/gregdel/homed/lib/components/common"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	components.Register(components.TypeESPHomeCover, NewCover)
}

// Cover is a component that controls a cover
type Cover struct {
	common.GenericSensor
}

// NewCover returns a new cover component
func NewCover() components.Component {
	return &Cover{}
}

// Type implements the Component interface
func (c *Cover) Type() components.Type {
	return components.TypeESPHomeCover
}

// Collectors implements the Component interface
func (c *Cover) Collectors(labels prometheus.Labels) []prometheus.Collector {
	return nil
}
