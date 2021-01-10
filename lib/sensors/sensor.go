package sensors

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

// Type reprensents a sensor type
type Type string

// Types
var (
	TypeUnknown         Type = "unknown"
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
	Collectors(labels prometheus.Labels) []prometheus.Collector
}

var registeredSensors map[string]func() Sensor

func register(t Type, fn func() Sensor) {
	if registeredSensors == nil {
		registeredSensors = map[string]func() Sensor{}
	}

	name := string(t)
	_, ok := registeredSensors[name]
	if ok {
		err := fmt.Errorf("sensors: sensor %s already regitered", name)
		panic(err)
	}

	registeredSensors[name] = fn
}

// New returns a new sensor from a type name
func New(typeName string) (Sensor, error) {
	fn, ok := registeredSensors[typeName]
	if !ok {
		return nil, fmt.Errorf("sensors: sensor %s is not registered", typeName)
	}

	return fn(), nil
}
