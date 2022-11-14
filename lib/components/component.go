package components

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gregdel/homed/lib/config"
	"github.com/prometheus/client_golang/prometheus"
)

// Type reprensents a component type
type Type string

// Types
var (
	// Common components
	TypePowerMeter    Type = "power_meter"
	TypeCounter       Type = "counter"
	TypeSwitch        Type = "switch"
	TypeBinaryLight   Type = "binary_light"
	TypeWifiSignal    Type = "wifi_signal"
	TypeDeviceStatus  Type = "device_status"
	TypeRollerShutter Type = "roller_shutter"

	// Other components
	TypeHomedTemperature    Type = "homed_temperature"
	TypeZigbeeTRV           Type = "zigbee_trv"
	TypeBoiler              Type = "boiler"
	TypeZigbeeClimateSensor Type = "zigbee_climate_sensor"
	TypeLinky               Type = "linky"
)

// Component represents a component
type Component interface {
	// Type represents the component type
	Type() Type

	// FriendlyName
	FriendlyName() string
	SetFriendlyName(string)

	// Device
	Device() *Device
	SetDevice(d *Device)

	// Config
	Config() *config.Component
	SetConfig(*config.Component)

	// Update is called every time a new mqtt payload is received on the state
	// topic
	Update([]byte) error

	// PostUpdate is a function called after the Update function
	PostUpdate() error

	// Collector returns the prometheus collectors to register for the component
	Collectors(labels prometheus.Labels) []prometheus.Collector

	// ReadOnly tells if the component is readonly
	ReadOnly() bool
	// Internal() tells if the component is internal
	Internal() bool

	// WriteCommand writes a command to a mqtt topic (for external components)
	WriteCommand(data []byte) error
	// ExecCommand executes a command (for internal components)
	ExecCommand(data []byte) error
	// PublishToStateTopic publishes raw data to the state topic
	PublishToStateTopic(data []byte) error

	SetStateTopic(string)
	SetCommandTopic(string)
	SetInternal(bool)
	SetMQTTClient(mqtt.Client)
	MQTTClient() mqtt.Client

	// ID
	ID() string
	SetID(id string)
}
