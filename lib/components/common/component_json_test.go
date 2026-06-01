package common

import (
	"encoding/json"
	"testing"

	"github.com/gregdel/homed/lib/components"
	"github.com/gregdel/homed/lib/config"
)

func TestComponentJSONUsesBinarySensorSnapshot(t *testing.T) {
	sensor := NewBinarySensor().(*BinarySensor)
	sensor.SetID("door")
	sensor.SetConfig(&config.Component{FriendlyName: "Door"})
	sensor.SetDevice(components.NewDevice("contact", "entry"))
	sensor.SetOn(true)

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

func TestComponentJSONUsesGenericSensorSnapshot(t *testing.T) {
	sensor := NewGenericSensor().(*GenericSensor)
	sensor.SetID("temperature")
	sensor.SetConfig(&config.Component{FriendlyName: "Temperature"})
	sensor.SetDevice(components.NewDevice("meter", "office"))
	sensor.SetSensorValue(21.5)

	data, err := json.Marshal(components.NewComponentJSON(sensor, true))
	if err != nil {
		t.Fatalf("failed to marshal component JSON: %s", err)
	}

	var got struct {
		Type   string `json:"type"`
		Values struct {
			ID    string  `json:"id"`
			Value float64 `json:"value"`
		} `json:"values"`
		HasGraph bool `json:"has_graph"`
	}
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("failed to decode component JSON: %s", err)
	}

	if got.Type != string(components.TypeGenericSensor) {
		t.Fatalf("unexpected type: got %q, want %q", got.Type, components.TypeGenericSensor)
	}
	if got.Values.ID != "temperature" {
		t.Fatalf("unexpected id: got %q, want %q", got.Values.ID, "temperature")
	}
	if got.Values.Value != 21.5 {
		t.Fatalf("unexpected value: got %f, want %f", got.Values.Value, 21.5)
	}
	if !got.HasGraph {
		t.Fatal("expected has_graph to be true")
	}
}

func TestDeviceStatusUpdateAndSnapshot(t *testing.T) {
	status := NewDeviceStatus().(*DeviceStatus)
	status.SetID("status")
	status.SetDevice(components.NewDevice("router", "hall"))

	if err := status.Update([]byte("online")); err != nil {
		t.Fatalf("failed to update status: %s", err)
	}

	if !status.IsOnline() {
		t.Fatal("expected status to be online")
	}
	if !status.Device().IsOnline() {
		t.Fatal("expected device to be online")
	}

	data, err := json.Marshal(components.NewComponentJSON(status, false))
	if err != nil {
		t.Fatalf("failed to marshal component JSON: %s", err)
	}

	var got struct {
		Values struct {
			Online bool `json:"online"`
		} `json:"values"`
	}
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("failed to decode component JSON: %s", err)
	}
	if !got.Values.Online {
		t.Fatal("expected online to be true")
	}
}

func TestComponentJSONUsesVirtualSwitchSnapshot(t *testing.T) {
	sw := NewVirtualSwitch().(*VirtualSwitch)
	sw.SetID("all_lights")
	sw.SetDevice(components.NewDevice("virtual", "house"))
	sw.SetOn(true)
	sw.incrementCounter()

	data, err := json.Marshal(components.NewComponentJSON(sw, false))
	if err != nil {
		t.Fatalf("failed to marshal component JSON: %s", err)
	}

	var got struct {
		Type   string `json:"type"`
		Values struct {
			On      bool    `json:"on"`
			Counter float64 `json:"counter"`
		} `json:"values"`
	}
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("failed to decode component JSON: %s", err)
	}

	if got.Type != string(components.TypeVirtualSwitch) {
		t.Fatalf("unexpected type: got %q, want %q", got.Type, components.TypeVirtualSwitch)
	}
	if !got.Values.On {
		t.Fatal("expected on to be true")
	}
	if got.Values.Counter != 1 {
		t.Fatalf("unexpected counter: got %f, want %f", got.Values.Counter, 1.0)
	}
}
