package esphome

import (
	"github.com/gregdel/homed/lib/components"
	"github.com/gregdel/homed/lib/components/common"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	components.Register(components.TypeESPHomeSwitch, NewSwitch)
}

// Switch is a switch component
type Switch struct {
	common.Component
	common.Switch
}

// NewSwitch returns a new switch component
func NewSwitch() components.Component {
	return &Switch{
		Switch: common.NewSwitch([]byte("ON"), []byte("OFF")),
	}
}

// Type implements the Component interface
func (s *Switch) Type() components.Type {
	return components.TypeESPHomeSwitch
}

// Collectors implements the Component interface
func (s *Switch) Collectors(labels prometheus.Labels) []prometheus.Collector {
	return []prometheus.Collector{
		components.GaugeCollector("switch", labels,
			func() float64 {
				if s.On {
					return 1
				}
				return 0
			},
		),
	}
}
