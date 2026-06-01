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
	graphs map[string]bool

	schedules map[string]*schedule.Schedule
	devices   map[string]*Device
}

// New returns a new Components type
func New(dataPath string) *Components {
	return &Components{
		dataPath: dataPath,
		byID:     map[string]Component{},
		byRoom:   map[string][]string{},
		graphs:   map[string]bool{},

		schedules: map[string]*schedule.Schedule{},
		devices:   map[string]*Device{},
	}
}

// ComponentJSON represents the JSON structure of a Component
type ComponentJSON struct {
	Values   any    `json:"values"`
	Type     string `json:"type"`
	ReadOnly bool   `json:"read_only"`
	HasGraph bool   `json:"has_graph"`
}

// ValuesSnapshotter is implemented by components with an explicit HTTP and
// websocket values JSON snapshot.
type ValuesSnapshotter interface {
	ValuesSnapshot() any
}

// NewComponentJSON returns a ComponentJSON from a Component
func NewComponentJSON(c Component, hasGraph bool) *ComponentJSON {
	values := any(c)
	if snapshotter, ok := c.(ValuesSnapshotter); ok {
		values = snapshotter.ValuesSnapshot()
	}

	return &ComponentJSON{
		Values:   values,
		Type:     string(c.Type()),
		ReadOnly: c.ReadOnly(),
		HasGraph: hasGraph,
	}
}

// MarshalJSON implements the json.Marshaler interface
func (c *Components) MarshalJSON() ([]byte, error) {
	components := []*ComponentJSON{}
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, component := range c.byID {
		components = append(components, NewComponentJSON(component, c.graphs[component.ID()]))
	}

	return json.Marshal(components)
}

// HasGraph tells if a component has Prometheus metrics available for graphing.
func (c *Components) HasGraph(id string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.graphs[id]
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

	var id string
	if cfg.ID != "" {
		_, err := c.Get(cfg.ID)
		if err == nil {
			return nil, fmt.Errorf("component with id %q is already configured", cfg.ID)
		}
		id = cfg.ID
	} else {
		id = fmt.Sprintf("%s_%s", deviceName, component.Type())
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

	labels := prometheus.Labels{
		"friendly_name":  component.FriendlyName(),
		"component_type": string(component.Type()),
		"device":         deviceName,
		"room":           roomName,
		"id":             id,
	}
	collectors := component.Collectors(labels)

	c.mu.Lock()

	c.byID[id] = component
	c.graphs[id] = len(collectors) > 0

	if len(c.byRoom[roomName]) == 0 {
		c.byRoom[roomName] = []string{}
	}
	c.byRoom[roomName] = append(c.byRoom[roomName], id)

	c.mu.Unlock()

	for _, c := range collectors {
		if err := prometheus.Register(c); err != nil {
			return nil, err
		}
	}

	return component, nil
}

// TemperatureController returns the temperature controller of the room
func (c *Components) TemperatureController(room string) TemperatureControllerInternal {
	c.mu.Lock()
	ids, ok := c.byRoom[room]
	c.mu.Unlock()
	if !ok {
		return nil
	}

	for _, id := range ids {
		component, err := c.Get(id)
		if err != nil {
			// TODO: handle this
			continue
		}

		if component.Type() != TypeHomedTemperature {
			continue
		}

		controller, ok := component.(TemperatureControllerInternal)
		if !ok {
			return nil
		}

		return controller
	}

	return nil
}
