package mqttd

import (
	"log/slog"
	"testing"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gregdel/homed/lib/components"
	"github.com/gregdel/homed/lib/config"
	"github.com/prometheus/client_golang/prometheus"
)

func TestHandleCommandIgnoresRetainedMessages(t *testing.T) {
	component := &testCommandComponent{id: "script"}
	m := &mqttd{
		logger:    slog.Default(),
		cmdTopics: map[string]components.Component{"script/command": component},
	}

	m.handleCommand(nil, testMessage{
		topic:    "script/command",
		payload:  []byte("{}"),
		retained: true,
	})

	if component.execCount != 0 {
		t.Fatalf("expected no command executions, got %d", component.execCount)
	}
}

func TestHandleCommandExecutesLiveMessages(t *testing.T) {
	component := &testCommandComponent{id: "script"}
	m := &mqttd{
		logger:    slog.Default(),
		cmdTopics: map[string]components.Component{"script/command": component},
	}

	m.handleCommand(nil, testMessage{
		topic:   "script/command",
		payload: []byte("{}"),
	})

	if component.execCount != 1 {
		t.Fatalf("expected one command execution, got %d", component.execCount)
	}
	if string(component.payload) != "{}" {
		t.Fatalf("unexpected payload: got %q", string(component.payload))
	}
}

type testMessage struct {
	topic    string
	payload  []byte
	retained bool
}

func (m testMessage) Duplicate() bool {
	return false
}

func (m testMessage) Qos() byte {
	return 0
}

func (m testMessage) Retained() bool {
	return m.retained
}

func (m testMessage) Topic() string {
	return m.topic
}

func (m testMessage) MessageID() uint16 {
	return 0
}

func (m testMessage) Payload() []byte {
	return m.payload
}

func (m testMessage) Ack() {}

type testCommandComponent struct {
	id        string
	execCount int
	payload   []byte
}

func (c *testCommandComponent) Type() components.Type {
	return components.TypeScript
}

func (c *testCommandComponent) FriendlyName() string {
	return c.id
}

func (c *testCommandComponent) Device() *components.Device {
	return nil
}

func (c *testCommandComponent) SetDevice(*components.Device) {}

func (c *testCommandComponent) Config() *config.Component {
	return &config.Component{}
}

func (c *testCommandComponent) SetConfig(*config.Component) {}

func (c *testCommandComponent) LoggerWithFields(logger *slog.Logger) *slog.Logger {
	return logger
}

func (c *testCommandComponent) Update([]byte) error {
	return nil
}

func (c *testCommandComponent) PostUpdate() error {
	return nil
}

func (c *testCommandComponent) Collectors(prometheus.Labels) []prometheus.Collector {
	return nil
}

func (c *testCommandComponent) ReadOnly() bool {
	return false
}

func (c *testCommandComponent) Internal() bool {
	return true
}

func (c *testCommandComponent) WriteCommand([]byte) error {
	return nil
}

func (c *testCommandComponent) ExecCommand(data []byte) error {
	c.execCount++
	c.payload = append([]byte(nil), data...)
	return nil
}

func (c *testCommandComponent) PublishToStateTopic([]byte) error {
	return nil
}

func (c *testCommandComponent) SetMQTTClient(mqtt.Client) {}

func (c *testCommandComponent) MQTTClient() mqtt.Client {
	return nil
}

func (c *testCommandComponent) ID() string {
	return c.id
}

func (c *testCommandComponent) SetID(id string) {
	c.id = id
}

func (c *testCommandComponent) Subscribe(string, chan components.Event) {}

func (c *testCommandComponent) Notify() {}
