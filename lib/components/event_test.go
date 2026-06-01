package components

import (
	"testing"
	"time"
)

func TestNewEventChannelUsesConfiguredBufferSize(t *testing.T) {
	ch := NewEventChannel()
	if cap(ch) != EventChannelBufferSize {
		t.Fatalf("unexpected event channel capacity: got %d, want %d", cap(ch), EventChannelBufferSize)
	}
}

func TestEventHandlerNotifyDeliversToSubscriber(t *testing.T) {
	events := NewEventChannel()
	handler := &EventHandler{}
	handler.Subscribe("subscriber", events)

	handler.Notify(Event{ID: "component"})

	select {
	case event := <-events:
		if event.ID != "component" {
			t.Fatalf("unexpected event ID: got %q, want %q", event.ID, "component")
		}
	default:
		t.Fatal("expected event to be delivered")
	}
}

func TestEventHandlerNotifyDropsWhenSubscriberIsFull(t *testing.T) {
	events := make(chan Event, 1)
	events <- Event{ID: "queued"}

	handler := &EventHandler{}
	handler.Subscribe("subscriber", events)

	done := make(chan struct{})
	go func() {
		handler.Notify(Event{ID: "dropped"})
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Notify blocked on a full subscriber channel")
	}

	event := <-events
	if event.ID != "queued" {
		t.Fatalf("unexpected event ID: got %q, want %q", event.ID, "queued")
	}

	select {
	case event := <-events:
		t.Fatalf("unexpected dropped event delivered: %q", event.ID)
	default:
	}
}
