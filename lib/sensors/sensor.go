package sensors

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Type reprensents a sensor type
type Type string

// Types
var (
	TypeTemperature     Type = "temperature"
	TypeRTL433          Type = "rtl_433"
	TypeStatus          Type = "device_status"
	TypeHumidity        Type = "humidity"
	TypeWifiSignal      Type = "wifi_signal"
	TypeZigbee2MQTTTuya Type = "zigbee2mqtt_tuya"
)

// Sensor represents a sensor
type Sensor interface {
	Type() Type
	Update([]byte) error
	PostUpdate() error
	Collectors(labels prometheus.Labels) []prometheus.Collector
}
