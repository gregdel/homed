package sensor

import (
	"encoding/json"

	"github.com/gregdel/homed/lib/components"
	"github.com/gregdel/homed/lib/components/common"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/atomic"
)

func init() {
	components.Register(components.TypeZigbeeClimateSensor, NewSensor)
}

// Data represents the data of the sensor
type Data struct {
	Battery      atomic.Float64 `json:"battery"`
	HumidityV    atomic.Float64 `json:"humidity"`
	LinkQuality  atomic.Float64 `json:"linkquality"`
	Pressure     atomic.Float64 `json:"pressure"`
	TemperatureV atomic.Float64 `json:"temperature"`
	Voltage      atomic.Float64 `json:"voltage"`
}

// Sensor reprensents a zigbee climate sensor
type Sensor struct {
	common.Component
	Data
}

// NewSensor returns a new climate sensor
func NewSensor() components.Component {
	return &Sensor{}
}

// Type implements the Component interface
func (s *Sensor) Type() components.Type {
	return components.TypeZigbeeClimateSensor
}

// Update implements the Component interface
func (s *Sensor) Update(value []byte) error {
	return json.Unmarshal(value, &s.Data)
}

// Temperature implements the TemperatureGetter interface
func (s *Sensor) Temperature() (float64, error) {
	if !s.Device().IsOnline() {
		return 0, components.ErrDeviceOffline
	}

	return s.TemperatureV.Load(), nil
}

// Humidity implements the HumidityGetter interface
func (s *Sensor) Humidity() (float64, error) {
	if !s.Device().IsOnline() {
		return 0, components.ErrDeviceOffline
	}

	return s.HumidityV.Load(), nil
}

// Collectors implements the Component interface
func (s *Sensor) Collectors(labels prometheus.Labels) []prometheus.Collector {
	return []prometheus.Collector{
		components.GaugeCollector("temperature", labels,
			func() float64 { return s.TemperatureV.Load() },
		),
		components.GaugeCollector("pressure", labels,
			func() float64 { return s.Pressure.Load() },
		),
		components.GaugeCollector("humidity", labels,
			func() float64 { return s.HumidityV.Load() },
		),
	}
}
