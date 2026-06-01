package sensor

import (
	"encoding/json"
	"testing"

	"github.com/gregdel/homed/lib/components"
	"github.com/gregdel/homed/lib/config"
)

func TestSensorUpdateSnapshotAndAccessors(t *testing.T) {
	s := testSensor()

	err := s.Update([]byte(`{
		"battery": 99,
		"humidity": 45.5,
		"linkquality": 120,
		"pressure": 1012.5,
		"temperature": 19.5,
		"voltage": 3000
	}`))
	if err != nil {
		t.Fatalf("failed to update sensor: %s", err)
	}

	temperature, err := s.Temperature()
	if err != nil {
		t.Fatalf("failed to read temperature: %s", err)
	}
	if temperature != 19.5 {
		t.Fatalf("unexpected temperature: got %f, want %f", temperature, 19.5)
	}

	humidity, err := s.Humidity()
	if err != nil {
		t.Fatalf("failed to read humidity: %s", err)
	}
	if humidity != 45.5 {
		t.Fatalf("unexpected humidity: got %f, want %f", humidity, 45.5)
	}

	data, err := json.Marshal(components.NewComponentJSON(s, true))
	if err != nil {
		t.Fatalf("failed to marshal component JSON: %s", err)
	}

	got := decodeJSONMap(t, data)
	assertJSONValue(t, got, "type", string(components.TypeZigbeeClimateSensor))
	assertJSONValue(t, got, "read_only", true)
	assertJSONValue(t, got, "has_graph", true)

	values := decodeRawJSONMap(t, got["values"])
	assertJSONValue(t, values, "id", "sensor")
	assertJSONValue(t, values, "friendly_name", "Sensor")
	assertJSONValue(t, values, "battery", float64(99))
	assertJSONValue(t, values, "humidity", 45.5)
	assertJSONValue(t, values, "linkquality", float64(120))
	assertJSONValue(t, values, "pressure", 1012.5)
	assertJSONValue(t, values, "temperature", 19.5)
	assertJSONValue(t, values, "voltage", float64(3000))
}

func TestSensorUpdatePreservesOmittedFields(t *testing.T) {
	s := testSensor()

	if err := s.Update([]byte(`{"temperature": 19.5, "humidity": 45.5}`)); err != nil {
		t.Fatalf("failed to update sensor: %s", err)
	}
	if err := s.Update([]byte(`{"temperature": 20}`)); err != nil {
		t.Fatalf("failed to update sensor: %s", err)
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.TemperatureV != 20 {
		t.Fatalf("unexpected temperature: got %f, want %f", s.TemperatureV, 20.0)
	}
	if s.HumidityV != 45.5 {
		t.Fatalf("unexpected humidity: got %f, want %f", s.HumidityV, 45.5)
	}
}

func testSensor() *Sensor {
	s := NewSensor().(*Sensor)
	s.SetID("sensor")
	s.SetConfig(&config.Component{FriendlyName: "Sensor"})
	device := components.NewDevice("sensor", "living_room")
	device.SetOnline(true)
	s.SetDevice(device)
	return s
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
