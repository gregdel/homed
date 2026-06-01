package components

import (
	"encoding/json"
	"testing"
)

func TestDeviceMarshalIncludesOnlineState(t *testing.T) {
	device := NewDevice("bridge", "utility")
	device.SetOnline(true)

	data, err := json.Marshal(device)
	if err != nil {
		t.Fatalf("failed to marshal device: %s", err)
	}

	var got struct {
		Name   string `json:"name"`
		Room   string `json:"room"`
		Online bool   `json:"online"`
	}
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("failed to decode device: %s", err)
	}

	if got.Name != "bridge" {
		t.Fatalf("unexpected name: got %q, want %q", got.Name, "bridge")
	}
	if got.Room != "utility" {
		t.Fatalf("unexpected room: got %q, want %q", got.Room, "utility")
	}
	if !got.Online {
		t.Fatal("expected online to be true")
	}
}
