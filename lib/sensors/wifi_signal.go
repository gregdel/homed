package sensors

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

// WifiSignal is a sensor that handles temperatures
type WifiSignal struct {
	BaseSensor
	Value float64
}

// NewWifiSignal returns a new wifi signal sensor
func NewWifiSignal() *WifiSignal {
	return &WifiSignal{}
}

// Type implements the Sensor interface
func (s *WifiSignal) Type() Type {
	return TypeWifiSignal
}

// Collectors implements the Sensor interface
func (s *WifiSignal) Collectors(labels prometheus.Labels) []prometheus.Collector {
	return []prometheus.Collector{
		prometheus.NewGaugeFunc(
			prometheus.GaugeOpts{
				Name:        "homed_wifi_signal",
				ConstLabels: labels,
			},
			func() float64 { return s.Value },
		),
	}
}

// Update implements the Sensor interface
func (s *WifiSignal) Update(value []byte) error {
	v, err := strconv.ParseFloat(string(value), 64)
	if err != nil {
		return err
	}
	s.Value = v
	return nil
}
