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
		"color_temp": 370,
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
	assertJSONValue(t, values, "color_mode", "rgb")
	assertJSONValue(t, values, "color_temp", float64(370))
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
	l.colorMode = "cwww"
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
	assertJSONValue(t, payload, "color_mode", "cwww")

	color := decodeRawJSONMap(t, payload["color"])
	assertJSONValue(t, color, "c", float64(10))
	assertJSONValue(t, color, "w", float64(20))

	if _, ok := payload["color_temp"]; ok {
		t.Fatalf("unexpected color_temp field for cwww payload")
	}
}

func TestLightTurnOnUsesRGBPayload(t *testing.T) {
	l, client := testLight()

	l.mu.Lock()
	l.brightness = 42
	l.colorMode = "rgb"
	l.red = 10
	l.green = 20
	l.blue = 30
	l.mu.Unlock()

	if err := l.TurnOn(); err != nil {
		t.Fatalf("failed to turn on light: %s", err)
	}

	payload := decodeJSONMap(t, client.payload)
	assertJSONValue(t, payload, "state", "ON")
	assertJSONValue(t, payload, "brightness", float64(42))
	assertJSONValue(t, payload, "color_mode", "rgb")

	color := decodeRawJSONMap(t, payload["color"])
	assertJSONValue(t, color, "r", float64(10))
	assertJSONValue(t, color, "g", float64(20))
	assertJSONValue(t, color, "b", float64(30))

	if _, ok := payload["color_temp"]; ok {
		t.Fatalf("unexpected color_temp field for rgb payload")
	}
}

func TestLightTurnOnDefaultsUnknownModeToCWWW(t *testing.T) {
	l, client := testLight()

	l.mu.Lock()
	l.brightness = 42
	l.colorMode = "white"
	l.mu.Unlock()

	if err := l.TurnOn(); err != nil {
		t.Fatalf("failed to turn on light: %s", err)
	}

	payload := decodeJSONMap(t, client.payload)
	assertJSONValue(t, payload, "state", "ON")
	assertJSONValue(t, payload, "brightness", float64(42))
	assertJSONValue(t, payload, "color_mode", "cwww")

	color := decodeRawJSONMap(t, payload["color"])
	assertJSONValue(t, color, "c", float64(defaultColdWhite))
	assertJSONValue(t, color, "w", float64(defaultWarmWhite))

	if _, ok := payload["color_temp"]; ok {
		t.Fatalf("unexpected color_temp field for cwww payload")
	}
}

func TestLightUpdatePreservesOmittedFields(t *testing.T) {
	l, _ := testLight()

	if err := l.Update([]byte(`{
		"state": "ON",
		"brightness": 42,
		"color_mode": "rgb",
		"color_temp": 370,
		"color": {"r": 30, "g": 40, "b": 50}
	}`)); err != nil {
		t.Fatalf("failed to update light: %s", err)
	}

	if err := l.Update([]byte(`{"brightness": 100}`)); err != nil {
		t.Fatalf("failed to update partial light state: %s", err)
	}

	snapshot := l.ValuesSnapshot().(Snapshot)
	if !snapshot.On {
		t.Fatalf("expected omitted state to preserve on=true")
	}
	if snapshot.Brightness != 100 {
		t.Fatalf("unexpected brightness: got %d, want 100", snapshot.Brightness)
	}
	if snapshot.ColorMode != "rgb" {
		t.Fatalf("unexpected color mode: got %q, want rgb", snapshot.ColorMode)
	}
	if snapshot.ColorTemp != 370 {
		t.Fatalf("unexpected color temp: got %d, want 370", snapshot.ColorTemp)
	}
	if snapshot.Red != 30 || snapshot.Green != 40 || snapshot.Blue != 50 {
		t.Fatalf("unexpected rgb: got %d/%d/%d, want 30/40/50",
			snapshot.Red, snapshot.Green, snapshot.Blue)
	}
}

func TestLightUpdatePreservesOmittedColorChannels(t *testing.T) {
	l, _ := testLight()

	l.mu.Lock()
	l.coldWhite = 10
	l.warmWhite = 20
	l.red = 30
	l.green = 40
	l.blue = 50
	l.mu.Unlock()

	if err := l.Update([]byte(`{
		"color_mode": "rgb",
		"color": {"r": 60, "g": 70, "b": 80}
	}`)); err != nil {
		t.Fatalf("failed to update rgb light state: %s", err)
	}

	snapshot := l.ValuesSnapshot().(Snapshot)
	if snapshot.ColdWhite != 10 || snapshot.WarmWhite != 20 {
		t.Fatalf("unexpected white channels after rgb update: got %d/%d, want 10/20",
			snapshot.ColdWhite, snapshot.WarmWhite)
	}
	if snapshot.Red != 60 || snapshot.Green != 70 || snapshot.Blue != 80 {
		t.Fatalf("unexpected rgb channels after rgb update: got %d/%d/%d, want 60/70/80",
			snapshot.Red, snapshot.Green, snapshot.Blue)
	}

	if err := l.Update([]byte(`{
		"color_mode": "cwww",
		"color": {"c": 90, "w": 100}
	}`)); err != nil {
		t.Fatalf("failed to update cwww light state: %s", err)
	}

	snapshot = l.ValuesSnapshot().(Snapshot)
	if snapshot.ColdWhite != 90 || snapshot.WarmWhite != 100 {
		t.Fatalf("unexpected white channels after cwww update: got %d/%d, want 90/100",
			snapshot.ColdWhite, snapshot.WarmWhite)
	}
	if snapshot.Red != 60 || snapshot.Green != 70 || snapshot.Blue != 80 {
		t.Fatalf("unexpected rgb channels after cwww update: got %d/%d/%d, want 60/70/80",
			snapshot.Red, snapshot.Green, snapshot.Blue)
	}
}

func TestLightUpdateAllowsExplicitZeroColorChannels(t *testing.T) {
	l, _ := testLight()

	l.mu.Lock()
	l.coldWhite = 10
	l.warmWhite = 20
	l.red = 30
	l.green = 40
	l.blue = 50
	l.mu.Unlock()

	if err := l.Update([]byte(`{
		"color": {"c": 0, "w": 0, "r": 0, "g": 0, "b": 0}
	}`)); err != nil {
		t.Fatalf("failed to update zero color channels: %s", err)
	}

	snapshot := l.ValuesSnapshot().(Snapshot)
	if snapshot.ColdWhite != 0 || snapshot.WarmWhite != 0 {
		t.Fatalf("unexpected white channels: got %d/%d, want 0/0",
			snapshot.ColdWhite, snapshot.WarmWhite)
	}
	if snapshot.Red != 0 || snapshot.Green != 0 || snapshot.Blue != 0 {
		t.Fatalf("unexpected rgb channels: got %d/%d/%d, want 0/0/0",
			snapshot.Red, snapshot.Green, snapshot.Blue)
	}
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
