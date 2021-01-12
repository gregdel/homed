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

// MarshalJSON implements the json.Marshaler interface
func (s Components) MarshalJSON() ([]byte, error) {
	type componentWithType struct {
		Component `json:"values"`
		Type      string `json:"type"`
		ReadOnly  bool   `json:"read_only"`
	}

	components := make([]componentWithType, len(s))
	for i := 0; i < len(s); i++ {
		components[i] = componentWithType{
			Component: s[i],
			Type:      string(s[i].Type()),
			ReadOnly:  s[i].ReadOnly(),
		}
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

	component.SetCommandTopic(cfg.CommandTopic)

	uuid, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	component.setID(uuid)

	*s = append(*s, component)
	return component, nil
}
