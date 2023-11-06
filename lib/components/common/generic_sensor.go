package common

import (
	"math"
	"strconv"

	"github.com/gregdel/homed/lib/components"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/atomic"
)

func init() {
	components.Register(components.TypeGenericSensor, NewGenericSensor)
}

// NewGenericSensor returns a new generic sensor
func NewGenericSensor() components.Component {
	return &GenericSensor{}
}

// GenericSensor represents a generic sensor
type GenericSensor struct {
	Component

	Value atomic.Float64 `json:"value"`
}

// Type implements the Component interface
func (g *GenericSensor) Type() components.Type {
	return components.TypeGenericSensor
}

// Update implements the Component interface
func (g *GenericSensor) Update(value []byte) error {
	v, err := strconv.ParseFloat(string(value), 64)
	if err != nil {
		return err
	}

	// TODO: find a better solution to handle NaN values
	if math.IsNaN(v) {
		v = -99999
	}

	g.Value.Store(v)
	return nil
}

// SensorValue implements the Sensor interface
func (g *GenericSensor) SensorValue() float64 {
	return g.Value.Load()
}

// Collectors implements the Component interface
func (g *GenericSensor) Collectors(labels prometheus.Labels) []prometheus.Collector {
	return []prometheus.Collector{
		components.GaugeCollector("generic_sensor", labels,
			func() float64 { return g.SensorValue() },
		),
	}
}
