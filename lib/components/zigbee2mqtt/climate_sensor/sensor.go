package sensor

import (
	"encoding/json"
	"sync"

	"github.com/gregdel/homed/lib/components"
	"github.com/gregdel/homed/lib/components/common"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	components.Register(components.TypeZigbeeClimateSensor, NewSensor)
}

// Data represents the data of the sensor
type Data struct {
	Battery      float64 `json:"battery"`
	HumidityV    float64 `json:"humidity"`
	LinkQuality  float64 `json:"linkquality"`
	Pressure     float64 `json:"pressure"`
	TemperatureV float64 `json:"temperature"`
	Voltage      float64 `json:"voltage"`
}

type Snapshot struct {
	common.SnapshotBase
	Data
}

// Sensor reprensents a zigbee climate sensor
type Sensor struct {
	common.Component

	mu sync.RWMutex
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
	s.mu.Lock()
	defer s.mu.Unlock()

	data := s.Data
	if err := json.Unmarshal(value, &data); err != nil {
		return err
	}

	s.Data = data
	return nil
}

// Temperature implements the TemperatureGetter interface
func (s *Sensor) Temperature() (float64, error) {
	if !s.Device().IsOnline() {
		return 0, components.ErrDeviceOffline
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.TemperatureV, nil
}

// Humidity implements the HumidityGetter interface
func (s *Sensor) Humidity() (float64, error) {
	if !s.Device().IsOnline() {
		return 0, components.ErrDeviceOffline
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.HumidityV, nil
}

func (s *Sensor) ValuesSnapshot() any {
	s.mu.RLock()
	data := s.Data
	s.mu.RUnlock()

	return Snapshot{
		SnapshotBase: s.SnapshotBase(),
		Data:         data,
	}
}

func (s *Sensor) DataSnapshot() Data {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.Data
}

// Collectors implements the Component interface
func (s *Sensor) Collectors(labels prometheus.Labels) []prometheus.Collector {
	return []prometheus.Collector{
		components.GaugeCollector("temperature", labels,
			func() float64 {
				s.mu.RLock()
				defer s.mu.RUnlock()
				return s.TemperatureV
			},
		),
		components.GaugeCollector("pressure", labels,
			func() float64 {
				s.mu.RLock()
				defer s.mu.RUnlock()
				return s.Pressure
			},
		),
		components.GaugeCollector("humidity", labels,
			func() float64 {
				s.mu.RLock()
				defer s.mu.RUnlock()
				return s.HumidityV
			},
		),
	}
}
