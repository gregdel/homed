package components

import "sync"

// Event represents and event
type Event struct {
	ID string
}

// EventHandler is a generic event handler for use in components
type EventHandler struct {
	mu          sync.RWMutex
	Incoming    chan Event
	subscribers map[string]chan Event
}

// Subscribe is a function to handle subscribers
func (e *EventHandler) Subscribe(id string, ch chan Event) {
	if e.subscribers == nil {
		e.subscribers = map[string]chan Event{}
	}

	e.mu.Lock()
	e.subscribers[id] = ch
	e.mu.Unlock()
}

// Notify is a function to notify subscribers
func (e *EventHandler) Notify(event Event) {
	if e.subscribers == nil {
		return
	}

	for _, ch := range e.subscribers {
		ch <- event
	}
}
