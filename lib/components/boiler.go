package components

import (
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	register(TypeBoiler, NewBoiler)
}

// Boiler is a component that controls the boiler
type Boiler struct {
	baseComponent

	On bool `json:"on"`
}

// NewBoiler returns a new status component
func NewBoiler() Component {
	return &Boiler{}
}

// Type implements the Component interface
func (s *Boiler) Type() Type {
	return TypeBoiler
}

// Collectors implements the Component interface
func (s *Boiler) Collectors(labels prometheus.Labels) []prometheus.Collector {
	return []prometheus.Collector{
		prometheus.NewGaugeFunc(
			prometheus.GaugeOpts{
				Name:        "homed_boiler",
				ConstLabels: labels,
			},
			func() float64 {
				if s.On {
					return 1
				}
				return 0
			},
		),
	}
}

// Update implements the Component interface
func (s *Boiler) Update(value []byte) error {
	data := string(value)
	if data == "ON" {
		s.On = true
	} else {
		s.On = false
	}

	return nil
}
