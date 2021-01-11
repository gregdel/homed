package sensors

import (
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	register(TypeWifiSignal, NewWifiSignal)
}

// WifiSignal is a sensor that handles temperatures
type WifiSignal struct {
	float64Sensor
}

// NewWifiSignal returns a new wifi signal sensor
func NewWifiSignal() Sensor {
	return &WifiSignal{}
}

// Type implements the Sensor interface
func (s *WifiSignal) Type() Type {
	return TypeWifiSignal
}

// Collectors implements the Sensor interface
func (s *WifiSignal) Collectors(labels prometheus.Labels) []prometheus.Collector {
	return singleCollector(s, labels, func() float64 { return s.Value })
}
