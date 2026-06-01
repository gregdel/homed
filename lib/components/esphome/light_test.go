package esphome

import (
	"encoding/json"
	"testing"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gregdel/homed/lib/components"
	"github.com/gregdel/homed/lib/config"
)

func TestLightUpdateAndSnapshot(t *testing.T) {
	l, _ := testLight()

	err := l.Update([]byte(`{
		"state": "ON",
		"brightness": 42,
		"color_mode": "rgb",
		"color": {"c": 10, "w": 20, "r": 30, "g": 40, "b": 50}
	}`))
	if err != nil {
		t.Fatalf("failed to update light: %s", err)
	}

	data, err := json.Marshal(components.NewComponentJSON(l, true))
	if err != nil {
		t.Fatalf("failed to marshal component JSON: %s", err)
	}

	got := decodeJSONMap(t, data)
	assertJSONValue(t, got, "type", string(components.TypeEsphomeLight))
	assertJSONValue(t, got, "read_only", false)
	assertJSONValue(t, got, "has_graph", true)

	values := decodeRawJSONMap(t, got["values"])
	assertJSONValue(t, values, "id", "light")
	assertJSONValue(t, values, "friendly_name", "Light")
	assertJSONValue(t, values, "on", true)
	assertJSONValue(t, values, "brightness", float64(42))
	assertJSONValue(t, values, "color_mode", "")
	assertJSONValue(t, values, "cold_white", float64(10))
	assertJSONValue(t, values, "warm_white", float64(20))
	assertJSONValue(t, values, "red", float64(30))
	assertJSONValue(t, values, "green", float64(40))
	assertJSONValue(t, values, "blue", float64(50))
}

func TestLightTurnOnUsesLockedStateCopy(t *testing.T) {
	l, client := testLight()

	l.mu.Lock()
	l.brightness = 42
	l.colorMode = "white"
	l.coldWhite = 10
	l.warmWhite = 20
	l.mu.Unlock()

	if err := l.TurnOn(); err != nil {
		t.Fatalf("failed to turn on light: %s", err)
	}

	if client.publishCount != 1 {
		t.Fatalf("unexpected publish count: got %d, want %d", client.publishCount, 1)
	}

	payload := decodeJSONMap(t, client.payload)
	assertJSONValue(t, payload, "state", "ON")
	assertJSONValue(t, payload, "brightness", float64(42))
	assertJSONValue(t, payload, "color_mode", "white")

	color := decodeRawJSONMap(t, payload["color"])
	assertJSONValue(t, color, "c", float64(10))
	assertJSONValue(t, color, "w", float64(20))
}

func testLight() (*Light, *fakeMQTTClient) {
	client := &fakeMQTTClient{connected: true}
	l := NewLight().(*Light)
	l.SetID("light")
	l.SetConfig(&config.Component{
		FriendlyName: "Light",
		CommandTopic: "light/set",
	})
	device := components.NewDevice("light", "living_room")
	device.SetOnline(true)
	l.SetDevice(device)
	l.SetMQTTClient(client)
	return l, client
}

type fakeMQTTClient struct {
	connected    bool
	publishCount int
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

func (c *fakeMQTTClient) Publish(_ string, _ byte, _ bool, payload any) mqtt.Token {
	c.publishCount++
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
