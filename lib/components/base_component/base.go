package base

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

// Component represents a base component
type Component struct {
	CommandTopic string
	StateTopic   string
	IsInternal   bool

	UUID      uuid.UUID  `json:"uuid"`
	UpdatedAt *time.Time `json:"updated_at"`
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

// SetInternal implements the Component interface
func (c *Component) SetInternal(internal bool) {
	c.IsInternal = internal
}

// SetID implements the Component interface
func (c *Component) SetID(uuid uuid.UUID) {
	c.UUID = uuid
}

// ID implements the Component interface
func (c *Component) ID() uuid.UUID {
	return c.UUID
}

// Internal implements the Component interface
func (c *Component) Internal() bool {
	return c.IsInternal
}

// WriteCommand implements the Component interface
func (c *Component) WriteCommand(client mqtt.Client, data []byte) error {
	if client == nil {
		return fmt.Errorf("components: missing mqtt client")
	}

	if c.ReadOnly() {
		return fmt.Errorf("components: component is read only")
	}

	token := client.Publish(c.CommandTopic, 0, false, data)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

// ExecCommand implements the Component interface
func (c *Component) ExecCommand(client mqtt.Client, data []byte) error {
	if !c.IsInternal {
		return fmt.Errorf("components: only internal components have exec commands")
	}

	return nil
}
