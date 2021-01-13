package components

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
)

// Type reprensents a component type
type Type string

// Types
var (
	TypeHomedTemperature Type = "homed_temperature"
	TypeTemperature      Type = "temperature"
	TypeRTL433           Type = "rtl_433"
	TypeDeviceStatus     Type = "device_status"
	TypeHumidity         Type = "humidity"
	TypeWifiSignal       Type = "wifi_signal"
	TypeTuyaTRV          Type = "tuya_trv"
	TypeBoiler           Type = "boiler"
)

// Component represents a component
type Component interface {
	// Type represents the component type
	Type() Type
	// Update is called every time a new mqtt payload is received on the state
	// topic
	Update([]byte) error
	// PostUpdate is a function called after the Update function
	PostUpdate() error
	// Collector returns the prometheus collectors to register for the component
	Collectors(labels prometheus.Labels) []prometheus.Collector
	// ReadOnly tells if the component is readonly
	ReadOnly() bool
	// SetCommandTopic sets the command topic for the component
	SetCommandTopic(string)
	// WriteCommand writes a command to a mqtt topic
	WriteCommand(client mqtt.Client, data []byte) error
	// ID returns the component id
	ID() uuid.UUID

	// setId sets the id of a component
	setID(uuid uuid.UUID)
}
