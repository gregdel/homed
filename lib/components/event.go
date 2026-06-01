package components

import "sync"

// EventChannelBufferSize is the size of component event subscriber channels.
const EventChannelBufferSize = 16

// Event represents an event.
type Event struct {
	ID string
}

// EventHandler is a generic event handler for use in components
type EventHandler struct {
	mu          sync.RWMutex
	Incoming    chan Event
	subscribers map[string]chan Event
}

// NewEventChannel returns a bounded channel for component event subscribers.
func NewEventChannel() chan Event {
	return make(chan Event, EventChannelBufferSize)
}

// Subscribe is a function to handle subscribers
func (e *EventHandler) Subscribe(id string, ch chan Event) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.subscribers == nil {
		e.subscribers = map[string]chan Event{}
	}

	e.subscribers[id] = ch
}

// Notify is a function to notify subscribers
func (e *EventHandler) Notify(event Event) {
	e.mu.RLock()
	if e.subscribers == nil {
		e.mu.RUnlock()
		return
	}

	subscribers := make([]chan Event, 0, len(e.subscribers))
	for _, ch := range e.subscribers {
		subscribers = append(subscribers, ch)
	}
	e.mu.RUnlock()

	for _, ch := range subscribers {
		select {
		case ch <- event:
		default:
		}
	}
}
