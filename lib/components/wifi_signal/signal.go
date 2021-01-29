package wifisignal

import (
	"github.com/gregdel/homed/lib/components"
	"github.com/gregdel/homed/lib/components/common"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	components.Register(components.TypeWifiSignal, New)
}

// WifiSignal is a component that handles temperatures
type WifiSignal struct {
	common.Component
	common.GenericSensor
}

// New returns a new wifi signal component
func New() components.Component {
	return &WifiSignal{}
}

// Type implements the Component interface
func (s *WifiSignal) Type() components.Type {
	return components.TypeWifiSignal
}

// Collectors implements the Component interface
func (s *WifiSignal) Collectors(labels prometheus.Labels) []prometheus.Collector {
	return components.SingleCollector(s, labels, func() float64 { return s.Value })
}
