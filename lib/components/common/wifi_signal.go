package common

import (
	"github.com/gregdel/homed/lib/components"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	components.Register(components.TypeWifiSignal, NewWifiSignal)
}

// WifiSignal is a component that handles temperatures
type WifiSignal struct {
	GenericSensor
}

// NewWifiSignal returns a new wifi signal component
func NewWifiSignal() components.Component {
	return &WifiSignal{}
}

// Type implements the Component interface
func (s *WifiSignal) Type() components.Type {
	return components.TypeWifiSignal
}

// Collectors implements the Component interface
func (s *WifiSignal) Collectors(labels prometheus.Labels) []prometheus.Collector {
	return []prometheus.Collector{
		components.GaugeCollector("wifi_signal", labels,
			func() float64 { return s.SensorValue() },
		),
	}
}
