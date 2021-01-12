package components

import (
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	register(TypeWifiSignal, NewWifiSignal)
}

// WifiSignal is a component that handles temperatures
type WifiSignal struct {
	baseComponent
	float64Component
}

// NewWifiSignal returns a new wifi signal component
func NewWifiSignal() Component {
	return &WifiSignal{}
}

// Type implements the Component interface
func (s *WifiSignal) Type() Type {
	return TypeWifiSignal
}

// Collectors implements the Component interface
func (s *WifiSignal) Collectors(labels prometheus.Labels) []prometheus.Collector {
	return singleCollector(s, labels, func() float64 { return s.Value })
}
