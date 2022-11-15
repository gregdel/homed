package common

import (
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gregdel/homed/lib/components"
	"github.com/gregdel/homed/lib/config"
)

// Component represents a base component
type Component struct {
	CommandTopic string `json:"-"`
	StateTopic   string `json:"-"`
	IsInternal   bool   `json:"-"`
	Hide         bool   `json:"hide"`
	config       *config.Component

	Dev *components.Device `json:"device"`

	mqttClient mqtt.Client

	Cid       string     `json:"id"`
	UpdatedAt *time.Time `json:"updated_at"`
	Name      string     `json:"friendly_name"`
}

// PostUpdate implements the Component interface
func (c *Component) PostUpdate() error {
	now := time.Now()
	c.UpdatedAt = &now
	return nil
}

// ReadOnly implements the Component interface
func (c *Component) ReadOnly() bool {
	return c.CommandTopic == ""
}

// SetCommandTopic implements the Component interface
func (c *Component) SetCommandTopic(topic string) {
	c.CommandTopic = topic
}

// SetStateTopic implements the Component interface
func (c *Component) SetStateTopic(topic string) {
	c.StateTopic = topic
}

// SetMQTTClient implements the Component interface
func (c *Component) SetMQTTClient(client mqtt.Client) {
	c.mqttClient = client
}

// MQTTClient implements the Component interface
func (c *Component) MQTTClient() mqtt.Client {
	return c.mqttClient
}

// Config implements the Component interface
func (c *Component) Config() *config.Component {
	return c.config
}

// SetConfig implements the Component interface
func (c *Component) SetConfig(config *config.Component) {
	c.config = config
	c.Hide = config.Hide
	c.IsInternal = config.Internal
	c.CommandTopic = config.CommandTopic
	c.StateTopic = config.StateTopic
	c.Name = config.FriendlyName
}

// SetID implements the Component interface
func (c *Component) SetID(id string) {
	c.Cid = id
}

// ID implements the Component interface
func (c *Component) ID() string {
	return c.Cid
}

// Internal implements the Component interface
func (c *Component) Internal() bool {
	return c.IsInternal
}

// Device implements the Component interface
func (c *Component) Device() *components.Device {
	return c.Dev
}

// SetDevice implements the Component interface
func (c *Component) SetDevice(d *components.Device) {
	c.Dev = d
}

// FriendlyName implements the Component interface
func (c *Component) FriendlyName() string {
	return c.Name
}

// SetFriendlyName implements the Component interface
func (c *Component) SetFriendlyName(n string) {
	c.Name = n
}

// WriteCommand implements the Component interface
func (c *Component) WriteCommand(data []byte) error {
	if c.mqttClient == nil {
		return components.ErrMissingMQTTClient
	}

	if c.ReadOnly() {
		return components.ErrComponentReadOnly
	}

	if c.Device() == nil {
		return components.ErrMissingDevice
	}

	if !c.IsInternal && !c.Device().Online {
		return components.ErrDeviceOffline
	}

	token := c.mqttClient.Publish(c.CommandTopic, 0, false, data)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

// ExecCommand implements the Component interface
func (c *Component) ExecCommand(data []byte) error {
	if c.mqttClient == nil {
		return components.ErrMissingMQTTClient
	}

	if !c.IsInternal {
		return components.ErrExecNotInternal
	}

	token := c.mqttClient.Publish(c.StateTopic, 0, true, data)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

// PublishToStateTopic publishes data to the state topic
func (c *Component) PublishToStateTopic(data []byte) error {
	token := c.mqttClient.Publish(c.StateTopic, 0, true, data)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}
