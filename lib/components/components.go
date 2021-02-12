package components

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/prometheus/client_golang/prometheus"
)

// Components is a type that holds the components
type Components struct {
	mu sync.Mutex

	dataPath string

	byID     map[string]Component
	byRoom   map[string][]string
	byDevice map[string][]string
}

// New returns a new Components type
func New(dataPath string) *Components {
	return &Components{
		dataPath: dataPath,
		byID:     map[string]Component{},
		byRoom:   map[string][]string{},
		byDevice: map[string][]string{},
	}
}

// ComponentJSON represents the JSONn structure of a Component
type ComponentJSON struct {
	Component `json:"values"`
	Type      string `json:"type"`
	ReadOnly  bool   `json:"read_only"`
}

// NewComponentJSON returns a ComponentJSON from a Component
func NewComponentJSON(c Component) *ComponentJSON {
	return &ComponentJSON{
		Component: c,
		Type:      string(c.Type()),
		ReadOnly:  c.ReadOnly(),
	}
}

// MarshalJSON implements the json.Marshaler interface
func (c *Components) MarshalJSON() ([]byte, error) {
	components := []*ComponentJSON{}
	for _, component := range c.byID {
		components = append(components, NewComponentJSON(component))
	}

	return json.Marshal(components)
}

// Get gets a component by its id
func (c *Components) Get(id string) (Component, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	co, ok := c.byID[id]
	if !ok {
		return nil, ErrComponentNotFound
	}

	return co, nil
}

// List lists all the components
func (c *Components) List() []Component {
	c.mu.Lock()
	defer c.mu.Unlock()

	components := []Component{}
	for _, component := range c.byID {
		components = append(components, component)
	}

	return components
}

// ListByRoom lists all the components by room
func (c *Components) ListByRoom(room string) []Component {
	components := []Component{}

	c.mu.Lock()
	ids, ok := c.byRoom[room]
	c.mu.Unlock()
	if !ok {
		return components
	}

	for _, id := range ids {
		component, err := c.Get(id)
		if err != nil {
			// TODO: handle this
			continue
		}
		components = append(components, component)
	}

	return components
}

// Add adds a component to the component slice
func (c *Components) Add(cfg Config, client mqtt.Client, roomName, deviceName string) (Component, error) {
	component, err := newComponent(cfg.Type)
	if err != nil {
		return nil, err
	}

	id := fmt.Sprintf("%s_%s", deviceName, component.Type())
	i := 1
	for {
		newID := fmt.Sprintf("%s_%d", id, i)
		_, err := c.Get(newID)
		if err == nil {
			i++
			continue
		}

		if err == ErrComponentNotFound {
			id = newID
			break
		}

		return nil, err
	}

	component.SetID(id)

	labels := prometheus.Labels{
		"device": deviceName,
		"room":   roomName,
		"id":     id,
	}

	collectors := component.Collectors(labels)
	for _, c := range collectors {
		if err := prometheus.Register(c); err != nil {
			return nil, err
		}

	}

	// TODO: find a better solution
	component.SetCommandTopic(cfg.CommandTopic)
	component.SetStateTopic(cfg.StateTopic)
	component.SetInternal(cfg.Internal)
	component.SetMQTTClient(client)
	component.SetRoom(roomName)
	component.SetDevice(deviceName)
	component.SetFriendlyName(cfg.FriendlyName)

	if sc, ok := component.(Scheduled); ok {
		// TODO: check and log the error
		sc.LoadSchedule(c.SchedulePath(component))
	}

	c.mu.Lock()

	c.byID[id] = component

	if len(c.byRoom[roomName]) == 0 {
		c.byRoom[roomName] = []string{}
	}
	c.byRoom[roomName] = append(c.byRoom[roomName], id)

	if len(c.byDevice[deviceName]) == 0 {
		c.byDevice[deviceName] = []string{}
	}
	c.byDevice[deviceName] = append(c.byDevice[deviceName], id)

	c.mu.Unlock()

	return component, nil
}

// SchedulePath returns the schedule file path of a component
func (c *Components) SchedulePath(component Component) string {
	return filepath.Join(c.dataPath, component.ID()+".yaml")
}
