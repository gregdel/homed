package homedtemperature

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/gregdel/homed/lib/components"
	"github.com/gregdel/homed/lib/config"
)

func TestStateSnapshotJSONShape(t *testing.T) {
	h := New().(*HomedTemperature)
	h.setStateForTest(Data{
		Current:       19.5,
		Target:        20.5,
		Mode:          string(components.TemperatureModeFixed),
		ManualTarget:  21,
		Heating:       true,
		Opportunistic: false,
		On:            true,
	})

	data, err := json.Marshal(h.StateSnapshot())
	if err != nil {
		t.Fatalf("failed to marshal state snapshot: %s", err)
	}

	got := decodeJSONMap(t, data)
	assertJSONKeys(t, got,
		"current",
		"target",
		"mode",
		"manual_target",
		"manual_until",
		"heating",
		"opportunistic",
		"on",
	)
	assertJSONValue(t, got, "current", 19.5)
	assertJSONValue(t, got, "target", 20.5)
	assertJSONValue(t, got, "mode", string(components.TemperatureModeFixed))
	assertJSONValue(t, got, "manual_target", float64(21))
	assertJSONValue(t, got, "manual_until", nil)
	assertJSONValue(t, got, "heating", true)
	assertJSONValue(t, got, "opportunistic", false)
	assertJSONValue(t, got, "on", true)
}

func TestComponentJSONUsesHomedTemperatureSnapshot(t *testing.T) {
	h := New().(*HomedTemperature)
	updatedAt := time.Date(2026, 6, 1, 12, 30, 0, 0, time.UTC)
	manualUntil := time.Date(2026, 6, 1, 14, 0, 0, 0, time.UTC)

	h.SetID("living_room_temperature")
	h.SetConfig(&config.Component{
		FriendlyName: "Living Room",
		Hide:         true,
	})
	h.SetDevice(components.NewDevice("thermostat", "living_room"))
	h.UpdatedAt.Store(&updatedAt)
	h.setStateForTest(Data{
		Current:       19.5,
		Target:        20.5,
		Mode:          string(components.TemperatureModeUntilDate),
		ManualTarget:  21,
		ManualUntil:   &manualUntil,
		Heating:       true,
		Opportunistic: false,
		On:            true,
	})

	data, err := json.Marshal(components.NewComponentJSON(h, true))
	if err != nil {
		t.Fatalf("failed to marshal component JSON: %s", err)
	}

	envelope := decodeJSONMap(t, data)
	assertJSONValue(t, envelope, "type", string(components.TypeHomedTemperature))
	assertJSONValue(t, envelope, "read_only", true)
	assertJSONValue(t, envelope, "has_graph", true)

	values := decodeRawJSONMap(t, envelope["values"])
	assertJSONKeys(t, values,
		"id",
		"updated_at",
		"friendly_name",
		"hide",
		"device",
		"current",
		"target",
		"mode",
		"manual_target",
		"manual_until",
		"heating",
		"opportunistic",
		"on",
	)
	assertJSONValue(t, values, "id", "living_room_temperature")
	assertJSONValue(t, values, "updated_at", updatedAt.Format(time.RFC3339))
	assertJSONValue(t, values, "friendly_name", "Living Room")
	assertJSONValue(t, values, "hide", true)
	assertJSONValue(t, values, "current", 19.5)
	assertJSONValue(t, values, "target", 20.5)
	assertJSONValue(t, values, "mode", string(components.TemperatureModeUntilDate))
	assertJSONValue(t, values, "manual_target", float64(21))
	assertJSONValue(t, values, "manual_until", manualUntil.Format(time.RFC3339))
	assertJSONValue(t, values, "heating", true)
	assertJSONValue(t, values, "opportunistic", false)
	assertJSONValue(t, values, "on", true)

	device := decodeRawJSONMap(t, values["device"])
	assertJSONValue(t, device, "name", "thermostat")
	assertJSONValue(t, device, "room", "living_room")
}

func TestUpdateReplacesStateFromRetainedPayload(t *testing.T) {
	h := New().(*HomedTemperature)
	manualUntil := time.Date(2026, 6, 1, 14, 0, 0, 0, time.UTC)

	payload := struct {
		Current       float64    `json:"current"`
		Target        float64    `json:"target"`
		Mode          string     `json:"mode"`
		ManualTarget  float64    `json:"manual_target"`
		ManualUntil   *time.Time `json:"manual_until"`
		Heating       bool       `json:"heating"`
		Opportunistic bool       `json:"opportunistic"`
		On            bool       `json:"on"`
	}{
		Current:       18.5,
		Target:        19,
		Mode:          string(components.TemperatureModeUntilDate),
		ManualTarget:  21,
		ManualUntil:   &manualUntil,
		Heating:       true,
		Opportunistic: false,
		On:            true,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal payload: %s", err)
	}

	if err := h.Update(data); err != nil {
		t.Fatalf("failed to update state: %s", err)
	}

	got := h.StateSnapshot()
	if got.Current != 18.5 ||
		got.Target != 19 ||
		got.Mode != string(components.TemperatureModeUntilDate) ||
		got.ManualTarget != 21 ||
		got.ManualUntil == nil ||
		!got.ManualUntil.Equal(manualUntil) ||
		!got.Heating ||
		got.Opportunistic ||
		!got.On {
		t.Fatalf("unexpected state snapshot: %#v", got)
	}
}

func TestUpdatePreservesOmittedFields(t *testing.T) {
	h := New().(*HomedTemperature)
	manualUntil := time.Date(2026, 6, 1, 14, 0, 0, 0, time.UTC)
	h.setStateForTest(Data{
		Current:       18.5,
		Target:        19,
		Mode:          string(components.TemperatureModeUntilDate),
		ManualTarget:  21,
		ManualUntil:   &manualUntil,
		Heating:       true,
		Opportunistic: true,
		On:            true,
	})

	if err := h.Update([]byte(`{"current":20}`)); err != nil {
		t.Fatalf("failed to update state: %s", err)
	}

	got := h.StateSnapshot()
	if got.Current != 20 ||
		got.Target != 19 ||
		got.Mode != string(components.TemperatureModeUntilDate) ||
		got.ManualTarget != 21 ||
		got.ManualUntil == nil ||
		!got.ManualUntil.Equal(manualUntil) ||
		!got.Heating ||
		!got.Opportunistic ||
		!got.On {
		t.Fatalf("unexpected state snapshot: %#v", got)
	}
}

func TestLockedAccessorsReadState(t *testing.T) {
	h := New().(*HomedTemperature)
	h.setStateForTest(Data{
		Target:       19,
		Mode:         string(components.TemperatureModeAuto),
		ManualTarget: 21,
		Heating:      true,
		On:           true,
	})

	target, err := h.TemperatureTarget()
	if err != nil {
		t.Fatalf("unexpected target error: %s", err)
	}
	if target != 19 {
		t.Fatalf("unexpected auto target: got %f, want %f", target, 19.0)
	}
	if !h.IsHeating() {
		t.Fatal("expected heating to be true")
	}
	if !h.IsOn() {
		t.Fatal("expected on to be true")
	}

	h.setStateForTest(Data{
		Target:       19,
		Mode:         string(components.TemperatureModeFixed),
		ManualTarget: 21,
		Heating:      true,
		On:           true,
	})

	target, err = h.TemperatureTarget()
	if err != nil {
		t.Fatalf("unexpected target error: %s", err)
	}
	if target != 21 {
		t.Fatalf("unexpected manual target: got %f, want %f", target, 21.0)
	}
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

func assertJSONKeys(t *testing.T, got map[string]json.RawMessage, keys ...string) {
	t.Helper()

	if len(got) != len(keys) {
		t.Fatalf("unexpected key count: got %d, want %d; object: %s", len(got), len(keys), mustMarshalJSON(t, got))
	}

	for _, key := range keys {
		if _, ok := got[key]; !ok {
			t.Fatalf("missing JSON key %q in object: %s", key, mustMarshalJSON(t, got))
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

func mustMarshalJSON(t *testing.T, value any) string {
	t.Helper()

	data, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("failed to marshal JSON: %s", err)
	}

	return string(data)
}

func (h *HomedTemperature) setStateForTest(data Data) {
	h.mu.Lock()
	defer h.mu.Unlock()

	data.ManualUntil = copyTimePtr(data.ManualUntil)
	h.Data = data
}
