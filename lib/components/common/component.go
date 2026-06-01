package common

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gregdel/homed/lib/components"
	"github.com/gregdel/homed/lib/config"
	"gopkg.in/yaml.v3"
)

// Timeout to wait for MQTT publish method
const publishTimeout = 15 * time.Second

// Component represents a base component
type Component struct {
	mu sync.RWMutex

	CommandTopic string `json:"-"`
	StateTopic   string `json:"-"`
	IsInternal   bool   `json:"-"`
	Hide         bool   `json:"hide"`
	config       *config.Component
	Events       components.EventHandler `json:"-"`

	Dev *components.Device `json:"device"`

	mqttClient mqtt.Client

	Cid        string    `json:"id"`
	UpdatedAt  Time      `json:"updated_at"`
	Name       string    `json:"friendly_name"`
	YAMLParams yaml.Node `json:"-"`
}

// SnapshotBase represents the common HTTP and websocket JSON fields shared by
// all components.
type SnapshotBase struct {
	ID           string             `json:"id"`
	UpdatedAt    *time.Time         `json:"updated_at"`
	FriendlyName string             `json:"friendly_name"`
	Hide         bool               `json:"hide"`
	Device       *components.Device `json:"device"`
}

// SnapshotBase returns a concurrency-safe copy of the common component fields.
func (c *Component) SnapshotBase() SnapshotBase {
	c.mu.RLock()
	id := c.Cid
	friendlyName := c.Name
	hide := c.Hide
	device := c.Dev
	c.mu.RUnlock()

	return SnapshotBase{
		ID:           id,
		UpdatedAt:    c.UpdatedAt.Load(),
		FriendlyName: friendlyName,
		Hide:         hide,
		Device:       device,
	}
}

// ValuesSnapshot returns the default HTTP and websocket JSON values shape.
func (c *Component) ValuesSnapshot() any {
	return c.SnapshotBase()
}

// PostUpdate implements the Component interface
func (c *Component) PostUpdate() error {
	now := time.Now()
	c.UpdatedAt.Store(&now)
	c.Events.Notify(components.Event{ID: c.ID()})
	return nil
}

// ReadOnly implements the Component interface
func (c *Component) ReadOnly() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.CommandTopic == ""
}

// LoggerWithFields implements the Component interface
func (c *Component) LoggerWithFields(logger *slog.Logger) *slog.Logger {
	log := logger.With(slog.String("device", c.Device().Name))

	if c.FriendlyName() != "" {
		log = log.With(slog.String("friendly_name", c.FriendlyName()))
	} else {
		log = log.With(slog.String("id", c.ID()))
	}

	return log
}

// SetMQTTClient implements the Component interface
func (c *Component) SetMQTTClient(client mqtt.Client) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.mqttClient = client
}

// MQTTClient implements the Component interface
func (c *Component) MQTTClient() mqtt.Client {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.mqttClient
}

// Config implements the Component interface
func (c *Component) Config() *config.Component {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.config
}

// SetConfig implements the Component interface
func (c *Component) SetConfig(config *config.Component) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.config = config
	c.Hide = config.Hide
	c.IsInternal = config.Internal
	c.CommandTopic = config.CommandTopic
	c.StateTopic = config.StateTopic
	c.Name = config.FriendlyName
	c.YAMLParams = config.Params
}

// SetID implements the Component interface
func (c *Component) SetID(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Cid = id
}

// ID implements the Component interface
func (c *Component) ID() string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.Cid
}

// Internal implements the Component interface
func (c *Component) Internal() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.IsInternal
}

// Device implements the Component interface
func (c *Component) Device() *components.Device {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.Dev
}

// SetDevice implements the Component interface
func (c *Component) SetDevice(d *components.Device) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Dev = d
}

// FriendlyName implements the Component interface
func (c *Component) FriendlyName() string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.Name
}

// SetFriendlyName implements the Component interface
func (c *Component) SetFriendlyName(n string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Name = n
}

// WriteCommand implements the Component interface
func (c *Component) WriteCommand(data []byte) error {
	client := c.MQTTClient()
	if client == nil {
		return components.ErrMissingMQTTClient
	}

	if c.ReadOnly() {
		return components.ErrComponentReadOnly
	}

	device := c.Device()
	if device == nil {
		return components.ErrMissingDevice
	}

	if !c.Internal() && !device.IsOnline() {
		return components.ErrDeviceOffline
	}

	c.mu.RLock()
	commandTopic := c.CommandTopic
	c.mu.RUnlock()

	token := client.Publish(commandTopic, 0, false, data)
	if !token.WaitTimeout(publishTimeout) {
		return fmt.Errorf("timeout reached while publishing")
	}
	if token.Error() != nil {
		return token.Error()
	}

	return nil
}

// ExecCommand implements the Component interface
func (c *Component) ExecCommand(data []byte) error {
	if c.MQTTClient() == nil {
		return components.ErrMissingMQTTClient
	}

	if !c.Internal() {
		return components.ErrExecNotInternal
	}

	return c.PublishToStateTopic(data)
}

// PublishToStateTopic publishes data to the state topic
func (c *Component) PublishToStateTopic(data []byte) error {
	client := c.MQTTClient()
	if !client.IsConnected() || !client.IsConnectionOpen() {
		return components.ErrMQTTClientNotConnected
	}

	c.mu.RLock()
	stateTopic := c.StateTopic
	c.mu.RUnlock()

	token := client.Publish(stateTopic, 0, true, data)
	if !token.WaitTimeout(publishTimeout) {
		return fmt.Errorf("timeout reached while publishing")
	}
	if token.Error() != nil {
		return token.Error()
	}

	return nil
}

// Subscribe implements the Component interface
func (c *Component) Subscribe(id string, ch chan components.Event) {
	c.Events.Subscribe(id, ch)
}

// Notify implements the Component interface
func (c *Component) Notify() {
	event := components.Event{ID: c.ID()}
	c.Events.Notify(event)
}
