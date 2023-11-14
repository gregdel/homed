package components

import (
	"context"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gregdel/homed/lib/config"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

// Type reprensents a component type
type Type string

// Types
var (
	// Common components
	TypeGenericSensor Type = "generic_sensor"
	TypeBinarySensor  Type = "binary_sensor"
	TypePowerMeter    Type = "power_meter"
	TypeCounter       Type = "counter"
	TypeSwitch        Type = "switch"
	TypeBinaryLight   Type = "binary_light"
	TypeWifiSignal    Type = "wifi_signal"
	TypeDeviceStatus  Type = "device_status"
	TypeRollerShutter Type = "roller_shutter"

	// Other components
	TypeHomedTemperature    Type = "homed_temperature"
	TypeHomedHumidity       Type = "homed_humidity"
	TypeZigbeeTRV           Type = "zigbee_trv"
	TypeBoiler              Type = "boiler"
	TypeZigbeeClimateSensor Type = "zigbee_climate_sensor"
)

// Component represents a component
type Component interface {
	// Type represents the component type
	Type() Type

	// FriendlyName
	FriendlyName() string

	// Device
	Device() *Device
	SetDevice(d *Device)

	// Config
	Config() *config.Component
	SetConfig(*config.Component)

	// Logger
	LoggerWithFields(*zap.Logger) *zap.Logger

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

	SetMQTTClient(mqtt.Client)
	MQTTClient() mqtt.Client

	// Run runs a goroutine for a component
	Run(context.Context, *zap.Logger, *Components) error

	// ID
	ID() string
	SetID(id string)

	// Subscribe
	Subscribe(string, chan Event)
	Notify()
}
