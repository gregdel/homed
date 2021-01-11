package sensors

import (
	"encoding/json"

	"github.com/prometheus/client_golang/prometheus"
)

// Sensors is a type that holds the sensors
type Sensors []Sensor

// New returns a new Sensors type
func New() Sensors {
	return []Sensor{}
}

// MarshalJSON implements the json.Marshaler interface
func (s Sensors) MarshalJSON() ([]byte, error) {
	type sensorWithType struct {
		Sensor `json:"values"`
		Type   string `json:"type"`
	}

	sensors := make([]sensorWithType, len(s))
	for i := 0; i < len(s); i++ {
		sensors[i] = sensorWithType{
			Sensor: s[i],
			Type:   string(s[i].Type()),
		}
	}

	return json.Marshal(sensors)
}

// Add adds a sensor to the sensor slice
// TODO: check if we need to return the sensor
func (s *Sensors) Add(sensorType string, labels prometheus.Labels) (Sensor, error) {
	sensor, err := newSensor(sensorType)
	if err != nil {
		return nil, err
	}

	collectors := sensor.Collectors(labels)
	for _, c := range collectors {
		if err := prometheus.Register(c); err != nil {
			return nil, err
		}

	}

	*s = append(*s, sensor)
	return sensor, nil
}
