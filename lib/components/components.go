package components

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/gregdel/homed/lib/config"
	"github.com/gregdel/homed/lib/schedule"
	"github.com/prometheus/client_golang/prometheus"
)

// Components is a type that holds the components
type Components struct {
	mu sync.Mutex

	dataPath string

	byID   map[string]Component
	byRoom map[string][]string

	schedules map[string]*schedule.Schedule
	devices   map[string]*Device
}

// New returns a new Components type
func New(dataPath string) *Components {
	return &Components{
		dataPath: dataPath,
		byID:     map[string]Component{},
		byRoom:   map[string][]string{},

		schedules: map[string]*schedule.Schedule{},
		devices:   map[string]*Device{},
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
func (c *Components) Add(cfg config.Component, roomName, deviceName string) (Component, error) {
	device, ok := c.devices[deviceName]
	if !ok {
		device = NewDevice(deviceName, roomName)
	}
	c.devices[device.Name] = device

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

	if err := device.AddComponent(component); err != nil {
		return nil, err
	}

	// TODO: find a better solution
	component.SetDevice(device)
	component.SetConfig(&cfg)

	if sc, ok := component.(Scheduled); ok {
		if cfg.ScheduleName == "" {
			return nil, ErrMissingScheduleName
		}

		schedulePath := schedulePath(c.dataPath, cfg.ScheduleName)

		schedule, ok := c.schedules[cfg.ScheduleName]
		if !ok {
			schedule, err = loadSchedule(schedulePath)
			if err != nil {
				return nil, err
			}

			c.schedules[cfg.ScheduleName] = schedule
		}

		sc.SetSchedule(schedule, schedulePath, cfg.ScheduleName)
	}

	c.mu.Lock()

	c.byID[id] = component

	if len(c.byRoom[roomName]) == 0 {
		c.byRoom[roomName] = []string{}
	}
	c.byRoom[roomName] = append(c.byRoom[roomName], id)

	c.mu.Unlock()

	labels := prometheus.Labels{
		"friendly_name": component.FriendlyName(),
		"device":        deviceName,
		"room":          roomName,
		"id":            id,
	}

	collectors := component.Collectors(labels)
	for _, c := range collectors {
		if err := prometheus.Register(c); err != nil {
			return nil, err
		}
	}

	return component, nil
}
