package sensors

import (
	"encoding/json"

	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	register(TypeZigbee2MQTTTuya, NewZigbee2MQTTTuya)
}

// Zigbee2MQTTTuya is a sensor that handles temperatures
type Zigbee2MQTTTuya struct {
	BaseSensor `json:"-"`

	HeatingSetpoint float64 `json:"current_heating_setpoint"`
	Temperature     float64 `json:"local_temperature"`
	Position        float64 `json:"position"`
	BatteryLow      bool    `json:"battery_low"`
}

// NewZigbee2MQTTTuya returns a new wifi signal sensor
func NewZigbee2MQTTTuya() Sensor {
	return &Zigbee2MQTTTuya{}
}

// Type implements the Sensor interface
func (s *Zigbee2MQTTTuya) Type() Type {
	return TypeZigbee2MQTTTuya
}

// Collectors implements the Sensor interface
func (s *Zigbee2MQTTTuya) Collectors(labels prometheus.Labels) []prometheus.Collector {
	prefix := "homed_zigbee2mqtt_tuya_"
	return []prometheus.Collector{
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
}

// Update implements the Sensor interface
func (s *Zigbee2MQTTTuya) Update(value []byte) error {
	return json.Unmarshal(value, s)
}
