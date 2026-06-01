package homedhumidity

import (
	"encoding/json"
	"testing"

	"github.com/gregdel/homed/lib/components"
	"github.com/gregdel/homed/lib/config"
)

func TestComponentJSONUsesMetadataSnapshot(t *testing.T) {
	h := New().(*HomedHumidity)
	h.SetID("bathroom_humidity_controller")
	h.SetConfig(&config.Component{
		FriendlyName: "Bathroom humidity",
		Hide:         true,
	})
	h.SetDevice(components.NewDevice("humidity_controller", "bathroom"))
	h.Params = Params{
		Sensor:            "bathroom_sensor",
		Switch:            "bathroom_fan",
		HumidityThreshold: 65,
	}

	data, err := json.Marshal(components.NewComponentJSON(h, false))
	if err != nil {
		t.Fatalf("failed to marshal component JSON: %s", err)
	}

	var got struct {
		Type   string                     `json:"type"`
		Values map[string]json.RawMessage `json:"values"`
	}
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("failed to decode component JSON: %s", err)
	}

	if got.Type != string(components.TypeHomedHumidity) {
		t.Fatalf("unexpected type: got %q, want %q", got.Type, components.TypeHomedHumidity)
	}
	assertJSONKeys(t, got.Values, "id", "updated_at", "friendly_name", "hide", "device")
	assertJSONValue(t, got.Values, "id", "bathroom_humidity_controller")
	assertJSONValue(t, got.Values, "updated_at", nil)
	assertJSONValue(t, got.Values, "friendly_name", "Bathroom humidity")
	assertJSONValue(t, got.Values, "hide", true)

	device := decodeRawJSONMap(t, got.Values["device"])
	assertJSONValue(t, device, "name", "humidity_controller")
	assertJSONValue(t, device, "room", "bathroom")
}

func decodeRawJSONMap(t *testing.T, data json.RawMessage) map[string]json.RawMessage {
	t.Helper()

	var got map[string]json.RawMessage
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("failed to decode JSON object: %s", err)
	}

	return got
}

func assertJSONKeys(t *testing.T, got map[string]json.RawMessage, keys ...string) {
	t.Helper()

	if len(got) != len(keys) {
		t.Fatalf("unexpected key count: got %d, want %d", len(got), len(keys))
	}

	for _, key := range keys {
		if _, ok := got[key]; !ok {
			t.Fatalf("missing JSON key %q", key)
		}
	}
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
