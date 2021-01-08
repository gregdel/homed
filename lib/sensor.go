package homed

// SensorType reprensents a sensor type
type SensorType string

// SensorTypes
var (
	SensorTypeUnknown         SensorType = "unknown"
	SensorTypeTemperature     SensorType = "temperature"
	SensorTypeRTL433          SensorType = "rtl_433"
	SensorTypeStatus          SensorType = "device_status"
	SensorTypeHumidity        SensorType = "humidity"
	SensorTypeWifiSignal      SensorType = "wifi_signal"
	SensorTypeZigbee2MQTTTuya SensorType = "zigbee2mqtt_tuya"
)

// Sensor represents a sensor
type Sensor interface {
	MQTTTopic() string
	SetDevice(*Device)
	SetMQTTTopic(string)
	Type() SensorType
	Update([]byte) error
	Init() error
}

// BaseSensor represents a basic sensor
type BaseSensor struct {
	device    *Device
	mqttTopic string
}

// MQTTTopic implements the Sensor interface
func (s *BaseSensor) MQTTTopic() string {
	return s.mqttTopic
}

// SetMQTTTopic implements the Sensor interface
func (s *BaseSensor) SetMQTTTopic(value string) {
	s.mqttTopic = value
}

// Type implements the Sensor interface
func (s *BaseSensor) Type() SensorType {
	return SensorTypeUnknown
}

// Init implements the Sensor interface
func (s *BaseSensor) Init() error {
	return nil
}

// Update implements the Sensor interface
func (s *BaseSensor) Update(_ []byte) error {
	return nil
}

// SetDevice implements the Sensor interface
func (s *BaseSensor) SetDevice(d *Device) {
	s.device = d
}
