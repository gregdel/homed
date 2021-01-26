package components

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
)

// Components is a type that holds the components
type Components []Component

// New returns a new Components type
func New() Components {
	return []Component{}
}

// ComponentJSON represents the JSONn structure of a Component
type ComponentJSON struct {
	Component `json:"values"`
	Type      string `json:"type"`
	ReadOnly  bool   `json:"read_only"`
}

// NewComponentJSON returns a ComponentJSON from a Component
func NewComponentJSON(c Component) ComponentJSON {
	return ComponentJSON{
		Component: c,
		Type:      string(c.Type()),
		ReadOnly:  c.ReadOnly(),
	}
}

// MarshalJSON implements the json.Marshaler interface
func (s Components) MarshalJSON() ([]byte, error) {
	components := make([]ComponentJSON, len(s))
	for i := 0; i < len(s); i++ {
		components[i] = NewComponentJSON(s[i])
	}

	return json.Marshal(components)
}

// Add adds a component to the component slice
func (s *Components) Add(cfg Config, labels prometheus.Labels) (Component, error) {
	component, err := newComponent(cfg.Type)
	if err != nil {
		return nil, err
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

	uuid, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	component.SetID(uuid)

	*s = append(*s, component)
	return component, nil
}
