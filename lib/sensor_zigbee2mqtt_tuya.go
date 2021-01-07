package homed

import (
	"encoding/json"

	"github.com/prometheus/client_golang/prometheus"
)

// SensorZigbee2MQTTTuya is a sensor that handles temperatures
type SensorZigbee2MQTTTuya struct {
	BaseSensor `json:"-"`

	HeatingSetpoint float64 `json:"current_heating_setpoint"`
	Temperature     float64 `json:"local_temperature"`
	Position        float64 `json:"position"`
	BatteryLow      bool    `json:"battery_low"`
}

// NewSensorZigbee2MQTTTuya returns a new wifi signal sensor
func NewSensorZigbee2MQTTTuya() *SensorZigbee2MQTTTuya {
	return &SensorZigbee2MQTTTuya{}
}

// Type implements the Sensor interface
func (s *SensorZigbee2MQTTTuya) Type() SensorType {
	return SensorTypeZigbee2MQTTTuya
}

// Init implements the Sensor interface
func (s *SensorZigbee2MQTTTuya) Init() error {
	if s.device == nil {
		return ErrMissingDevice
	}

	if s.device.Room == nil {
		return ErrMissingRoom
	}

	labels := prometheus.Labels{
		"device": string(s.device.Name),
		"room":   string(s.device.Room.Name),
	}

	prefix := "homed_zigbee2mqtt_tuya_"

	collectors := []prometheus.Collector{
		prometheus.NewGaugeFunc(
			prometheus.GaugeOpts{
				Name:        prefix + "temperature",
				ConstLabels: labels,
			},
			func() float64 { return s.Temperature },
		),
		prometheus.NewGaugeFunc(
			prometheus.GaugeOpts{
				Name:        prefix + "heating_set_point",
				ConstLabels: labels,
			},
			func() float64 { return s.HeatingSetpoint },
		),
		prometheus.NewGaugeFunc(
			prometheus.GaugeOpts{
				Name:        prefix + "position",
				ConstLabels: labels,
			},
			func() float64 { return s.Position },
		),
		prometheus.NewGaugeFunc(
			prometheus.GaugeOpts{
				Name:        prefix + "battery_low",
				ConstLabels: labels,
			},
			func() float64 {
				if s.BatteryLow {
					return 1
				}
				return 0
			},
		),
	}

	for _, collector := range collectors {
		if err := prometheus.Register(collector); err != nil {
			return err
		}
	}

	return nil
}

// Update implements the Sensor interface
func (s *SensorZigbee2MQTTTuya) Update(value []byte) error {
	return json.Unmarshal(value, s)
}
