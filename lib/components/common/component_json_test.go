package common

import (
	"encoding/json"
	"testing"

	"github.com/gregdel/homed/lib/components"
	"github.com/gregdel/homed/lib/config"
)

func TestComponentJSONFallsBackToComponentValues(t *testing.T) {
	sensor := NewBinarySensor().(*BinarySensor)
	sensor.SetID("door")
	sensor.SetConfig(&config.Component{FriendlyName: "Door"})
	sensor.SetDevice(components.NewDevice("contact", "entry"))
	sensor.On.Store(true)

	data, err := json.Marshal(components.NewComponentJSON(sensor, false))
	if err != nil {
		t.Fatalf("failed to marshal component JSON: %s", err)
	}

	var got struct {
		Type   string `json:"type"`
		Values struct {
			ID           string `json:"id"`
			FriendlyName string `json:"friendly_name"`
			On           bool   `json:"on"`
		} `json:"values"`
		HasGraph bool `json:"has_graph"`
	}
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("failed to decode component JSON: %s", err)
	}

	if got.Type != string(components.TypeBinarySensor) {
		t.Fatalf("unexpected type: got %q, want %q", got.Type, components.TypeBinarySensor)
	}
	if got.Values.ID != "door" {
		t.Fatalf("unexpected id: got %q, want %q", got.Values.ID, "door")
	}
	if got.Values.FriendlyName != "Door" {
		t.Fatalf("unexpected friendly name: got %q, want %q", got.Values.FriendlyName, "Door")
	}
	if !got.Values.On {
		t.Fatal("expected on to be true")
	}
	if got.HasGraph {
		t.Fatal("expected has_graph to be false")
	}
}
