package common

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Component represents a base component
type Component struct {
	CommandTopic string `json:"-"`
	StateTopic   string `json:"-"`
	IsInternal   bool   `json:"-"`

	mqttClient mqtt.Client

	Cid        string     `json:"id"`
	UpdatedAt  *time.Time `json:"updated_at"`
	RoomName   string     `json:"room_name"`
	DeviceName string     `json:"device_name"`
}

// PostUpdate implements the Component interface
func (c *Component) PostUpdate() error {
	now := time.Now()
	c.UpdatedAt = &now
	return nil
}

// ReadOnly implements the Component interface
func (c *Component) ReadOnly() bool {
	if c.CommandTopic == "" {
		return true
	}

	return false
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

// SetInternal implements the Component interface
func (c *Component) SetInternal(internal bool) {
	c.IsInternal = internal
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

// Room implements the Component interface
func (c *Component) Room() string {
	return c.RoomName
}

// SetRoom implements the Component interface
func (c *Component) SetRoom(name string) {
	c.RoomName = name
}

// Device implements the Component interface
func (c *Component) Device() string {
	return c.DeviceName
}

// SetDevice implements the Component interface
func (c *Component) SetDevice(name string) {
	c.DeviceName = name
}

// WriteCommand implements the Component interface
func (c *Component) WriteCommand(data []byte) error {
	if c.mqttClient == nil {
		return fmt.Errorf("components: missing mqtt client")
	}

	if c.ReadOnly() {
		return fmt.Errorf("components: component is read only")
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
		return fmt.Errorf("components: missing mqtt client")
	}

	if !c.IsInternal {
		return fmt.Errorf("components: only internal components have exec commands")
	}

	token := c.mqttClient.Publish(c.StateTopic, 0, true, data)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}
