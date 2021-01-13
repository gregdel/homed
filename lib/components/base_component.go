package components

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

type baseComponent struct {
	commandTopic string
	stateTopic   string
	internal     bool

	UUID      uuid.UUID  `json:"uuid"`
	UpdatedAt *time.Time `json:"updated_at"`
}

// PostUpdate implements the Component interface
func (bs *baseComponent) PostUpdate() error {
	now := time.Now()
	bs.UpdatedAt = &now
	return nil
}

// ReadOnly implements the Component interface
func (bs *baseComponent) ReadOnly() bool {
	if bs.commandTopic == "" {
		return true
	}

	return false
}

// SetCommandTopic implements the Component interface
func (bs *baseComponent) SetCommandTopic(topic string) {
	bs.commandTopic = topic
}

// SetStateTopic implements the Component interface
func (bs *baseComponent) SetStateTopic(topic string) {
	bs.stateTopic = topic
}

// SetInternal implements the Component interface
func (bs *baseComponent) SetInternal(internal bool) {
	bs.internal = internal
}

// setID implements the Component interface
func (bs *baseComponent) setID(uuid uuid.UUID) {
	bs.UUID = uuid
}

// ID implements the Component interface
func (bs *baseComponent) ID() uuid.UUID {
	return bs.UUID
}

// Internal implements the Component interface
func (bs *baseComponent) Internal() bool {
	return bs.internal
}

// WriteCommand implements the Component interface
func (bs *baseComponent) WriteCommand(client mqtt.Client, data []byte) error {
	if client == nil {
		return fmt.Errorf("components: missing mqtt client")
	}

	if bs.ReadOnly() {
		return fmt.Errorf("components: component is read only")
	}

	token := client.Publish(bs.commandTopic, 0, false, data)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

// ExecCommand implements the Component interface
func (bs *baseComponent) ExecCommand(client mqtt.Client, data []byte) error {
	if !bs.internal {
		return fmt.Errorf("components: only internal components have exec commands")
	}

	return nil
}
