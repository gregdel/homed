package components

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type baseComponent struct {
	commandTopic string

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
