package boiler

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gregdel/homed/lib/components"
	"github.com/gregdel/homed/lib/config"
)

func TestWriteCommandFirstWriteSucceeds(t *testing.T) {
	b, client := testBoiler()

	if err := b.WriteCommand([]byte("ON")); err != nil {
		t.Fatalf("unexpected write error: %s", err)
	}

	if client.publishCount != 1 {
		t.Fatalf("unexpected publish count: got %d, want %d", client.publishCount, 1)
	}
	if string(client.payload) != "ON" {
		t.Fatalf("unexpected payload: got %q, want %q", client.payload, "ON")
	}
	if b.lastStateChangePtr() == nil {
		t.Fatal("expected last state change to be set")
	}
}

func TestWriteCommandRejectsInsideCooldown(t *testing.T) {
	b, client := testBoiler()

	b.mu.Lock()
	b.lastStateChange = time.Now()
	b.mu.Unlock()

	lastStateChange := b.lastStateChangePtr()
	err := b.reserveStateChange(lastStateChange.Add(cooldownDuration - time.Second))
	if !errors.Is(err, ErrBoilerCooldown) {
		t.Fatalf("unexpected cooldown error: got %v, want %v", err, ErrBoilerCooldown)
	}

	if err := b.WriteCommand([]byte("OFF")); !errors.Is(err, ErrBoilerCooldown) {
		t.Fatalf("unexpected write error: got %v, want %v", err, ErrBoilerCooldown)
	}
	if client.publishCount != 0 {
		t.Fatalf("unexpected publish count: got %d, want %d", client.publishCount, 0)
	}
}

func TestWriteCommandAllowsAfterCooldown(t *testing.T) {
	b, client := testBoiler()
	now := time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC)

	if err := b.reserveStateChange(now); err != nil {
		t.Fatalf("unexpected reserve error: %s", err)
	}
	if err := b.reserveStateChange(now.Add(cooldownDuration)); err != nil {
		t.Fatalf("unexpected reserve error at cooldown boundary: %s", err)
	}

	got := b.lastStateChangePtr()
	want := now.Add(cooldownDuration)
	if got == nil || !got.Equal(want) {
		t.Fatalf("unexpected last state change: got %v, want %v", got, want)
	}

	b.mu.Lock()
	b.lastStateChange = time.Now().Add(-cooldownDuration)
	b.mu.Unlock()

	if err := b.WriteCommand([]byte("OFF")); err != nil {
		t.Fatalf("unexpected write error: %s", err)
	}
	if client.publishCount != 1 {
		t.Fatalf("unexpected publish count: got %d, want %d", client.publishCount, 1)
	}
}

func TestBoilerSnapshotJSONShape(t *testing.T) {
	b, _ := testBoiler()

	data, err := json.Marshal(components.NewComponentJSON(b, true))
	if err != nil {
		t.Fatalf("failed to marshal component JSON: %s", err)
	}

	got := decodeJSONMap(t, data)
	assertJSONValue(t, got, "type", string(components.TypeBoiler))
	assertJSONValue(t, got, "read_only", false)
	assertJSONValue(t, got, "has_graph", true)

	values := decodeRawJSONMap(t, got["values"])
	assertJSONValue(t, values, "id", "boiler")
	assertJSONValue(t, values, "friendly_name", "Boiler")
	assertJSONValue(t, values, "hide", false)
	assertJSONValue(t, values, "on", false)
	assertJSONValue(t, values, "last_state_change", nil)

	now := time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC)
	if err := b.reserveStateChange(now); err != nil {
		t.Fatalf("unexpected reserve error: %s", err)
	}

	data, err = json.Marshal(components.NewComponentJSON(b, true))
	if err != nil {
		t.Fatalf("failed to marshal component JSON: %s", err)
	}
	got = decodeJSONMap(t, data)
	values = decodeRawJSONMap(t, got["values"])
	assertJSONValue(t, values, "last_state_change", now.Format(time.RFC3339))
}

func testBoiler() (*Boiler, *fakeMQTTClient) {
	client := &fakeMQTTClient{connected: true}
	b := NewBoiler().(*Boiler)
	b.SetID("boiler")
	b.SetConfig(&config.Component{
		FriendlyName: "Boiler",
		CommandTopic: "boiler/set",
	})
	device := components.NewDevice("boiler", "utility")
	device.Online.Store(true)
	b.SetDevice(device)
	b.SetMQTTClient(client)
	return b, client
}

type fakeMQTTClient struct {
	connected    bool
	publishCount int
	topic        string
	payload      []byte
}

func (c *fakeMQTTClient) IsConnected() bool {
	return c.connected
}

func (c *fakeMQTTClient) IsConnectionOpen() bool {
	return c.connected
}

func (c *fakeMQTTClient) Connect() mqtt.Token {
	return fakeToken{}
}

func (c *fakeMQTTClient) Disconnect(uint) {}

func (c *fakeMQTTClient) Publish(topic string, _ byte, _ bool, payload interface{}) mqtt.Token {
	c.publishCount++
	c.topic = topic
	switch p := payload.(type) {
	case []byte:
		c.payload = append([]byte(nil), p...)
	case string:
		c.payload = []byte(p)
	}
	return fakeToken{}
}

func (c *fakeMQTTClient) Subscribe(string, byte, mqtt.MessageHandler) mqtt.Token {
	return fakeToken{}
}

func (c *fakeMQTTClient) SubscribeMultiple(map[string]byte, mqtt.MessageHandler) mqtt.Token {
	return fakeToken{}
}

func (c *fakeMQTTClient) Unsubscribe(...string) mqtt.Token {
	return fakeToken{}
}

func (c *fakeMQTTClient) AddRoute(string, mqtt.MessageHandler) {}

func (c *fakeMQTTClient) OptionsReader() mqtt.ClientOptionsReader {
	return mqtt.ClientOptionsReader{}
}

type fakeToken struct{}

func (fakeToken) Wait() bool {
	return true
}

func (fakeToken) WaitTimeout(time.Duration) bool {
	return true
}

func (fakeToken) Done() <-chan struct{} {
	done := make(chan struct{})
	close(done)
	return done
}

func (fakeToken) Error() error {
	return nil
}

func decodeJSONMap(t *testing.T, data []byte) map[string]json.RawMessage {
	t.Helper()

	var got map[string]json.RawMessage
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("failed to decode JSON object: %s", err)
	}

	return got
}

func decodeRawJSONMap(t *testing.T, data json.RawMessage) map[string]json.RawMessage {
	t.Helper()
	return decodeJSONMap(t, data)
}

func assertJSONValue(t *testing.T, got map[string]json.RawMessage, key string, want any) {
	t.Helper()

	data, ok := got[key]
	if !ok {
		t.Fatalf("missing JSON key %q", key)
	}

	var value any
	if err := json.Unmarshal(data, &value); err != nil {
		t.Fatalf("failed to decode JSON key %q: %s", key, err)
	}

	if value != want {
		t.Fatalf("unexpected JSON value for %q: got %#v, want %#v", key, value, want)
	}
}
