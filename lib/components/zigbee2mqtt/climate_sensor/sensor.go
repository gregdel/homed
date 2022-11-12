package sensor

import (
	"encoding/json"

	"github.com/gregdel/homed/lib/components"
	"github.com/gregdel/homed/lib/components/common"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	components.Register(components.TypeZigbeeClimateSensor, NewSensor)
}

// Sensor reprensents a zigbee climate sensor
type Sensor struct {
	common.Component

	Battery     float64 `json:"battery"`
	Humidity    float64 `json:"humidity"`
	LinkQuality float64 `json:"linkquality"`
	Pressure    float64 `json:"pressure"`
	Temp        float64 `json:"temperature"`
	Voltage     float64 `json:"voltage"`
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
	return json.Unmarshal(value, s)
}

// Temperature implements the TemperatureGetter interface
func (s *Sensor) Temperature() (float64, error) {
	if !s.Device().Online {
		return 0, components.ErrDeviceOffline
	}

	return s.Temp, nil
}

// Collectors implements the Component interface
func (s *Sensor) Collectors(labels prometheus.Labels) []prometheus.Collector {
	return []prometheus.Collector{
		components.GaugeCollector("temperature", labels,
			func() float64 { return s.Temp },
		),
		components.GaugeCollector("pressure", labels,
			func() float64 { return s.Pressure },
		),
		components.GaugeCollector("humidity", labels,
			func() float64 { return s.Humidity },
		),
		components.GaugeCollector("link_quality", labels,
			func() float64 { return s.LinkQuality },
		),
	}
}
