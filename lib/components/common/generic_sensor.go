package common

import (
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/gregdel/homed/lib/components"
	"github.com/prometheus/client_golang/prometheus"
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

	mu    sync.RWMutex
	value float64
}

type GenericSensorSnapshot struct {
	ID           string             `json:"id"`
	UpdatedAt    *time.Time         `json:"updated_at"`
	FriendlyName string             `json:"friendly_name"`
	Hide         bool               `json:"hide"`
	Device       *components.Device `json:"device"`
	Value        float64            `json:"value"`
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

	g.SetSensorValue(v)
	return nil
}

// SensorValue implements the Sensor interface
func (g *GenericSensor) SensorValue() float64 {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.value
}

func (g *GenericSensor) SetSensorValue(value float64) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.value = value
}

func (g *GenericSensor) Snapshot() any {
	return GenericSensorSnapshot{
		ID:           g.ID(),
		UpdatedAt:    g.UpdatedAt.Load(),
		FriendlyName: g.FriendlyName(),
		Hide:         g.Hide,
		Device:       g.Device(),
		Value:        g.SensorValue(),
	}
}

// Collectors implements the Component interface
func (g *GenericSensor) Collectors(labels prometheus.Labels) []prometheus.Collector {
	return []prometheus.Collector{
		components.GaugeCollector("generic_sensor", labels,
			func() float64 { return g.SensorValue() },
		),
	}
}
