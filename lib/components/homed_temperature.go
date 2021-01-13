package components

import (
	"encoding/json"

	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	register(TypeHomedTemperature, NewHomedTemperature)
}

// HomedTemperatureData represents the data of HomedTemperature
type HomedTemperatureData struct {
	Current float64 `json:"current"`
	Target  float64 `json:"target"`
}

// HomedTemperature is a component that handles temperatures
type HomedTemperature struct {
	baseComponent
	HomedTemperatureData
}

// NewHomedTemperature returns a new temperature component
func NewHomedTemperature() Component {
	return &HomedTemperature{}
}

// Type implements the Component interface
func (s *HomedTemperature) Type() Type {
	return TypeHomedTemperature
}

// Collectors implements the Component interface
func (s *HomedTemperature) Collectors(labels prometheus.Labels) []prometheus.Collector {
	prefix := "homed_temperature_control_"
	return []prometheus.Collector{
		prometheus.NewGaugeFunc(
			prometheus.GaugeOpts{
				Name:        prefix + "current",
				ConstLabels: labels,
			},
			func() float64 { return s.Current },
		),
		prometheus.NewGaugeFunc(
			prometheus.GaugeOpts{
				Name:        prefix + "target",
				ConstLabels: labels,
			},
			func() float64 { return s.Target },
		),
	}
}

// Update implements the Component interface
func (s *HomedTemperature) Update(value []byte) error {
	return json.Unmarshal(value, s)
}
