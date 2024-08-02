package esphome

import (
	"encoding/json"

	"github.com/gregdel/homed/lib/components"
	"github.com/gregdel/homed/lib/components/common"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/atomic"
)

func init() {
	components.Register(components.TypeEsphomeLight, NewLight)
}

// Light represents a esphome light.
type Light struct {
	common.Component
	On atomic.Bool `json:"on"`

	Brightness atomic.Uint32 `json:"brightness"`
	ColorMode  atomic.String `json:"color_mode"`

	ColdWhite atomic.Uint32 `json:"cold_white"`
	WarmWhite atomic.Uint32 `json:"warm_white"`
	Red       atomic.Uint32 `json:"red"`
	Green     atomic.Uint32 `json:"green"`
	Blue      atomic.Uint32 `json:"blue"`
}

// NewLight returns a new light
func NewLight() components.Component {
	return &Light{}
}

type color struct {
	C uint8 `json:"c"`
	W uint8 `json:"w"`
	R uint8 `json:"r"`
	G uint8 `json:"g"`
	B uint8 `json:"b"`
}

type payload struct {
	State      string `json:"state"`
	Brightness uint8  `json:"brightness"`
	ColorMode  string `json:"color_mode"`
	Color      color  `json:"color"`
}

// TurnOn implements the switch interface
func (l *Light) updateState(s string) error {
	data := payload{
		State:      s,
		Brightness: uint8(l.Brightness.Load()),
		ColorMode:  l.ColorMode.Load(),
		Color: color{
			C: uint8(l.ColdWhite.Load()),
			W: uint8(l.WarmWhite.Load()),
		},
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return l.WriteCommand(payload)
}

// IsOn implements the BinarySensor interface
func (l *Light) IsOn() bool {
	return l.On.Load()
}

// TurnOn implements the switch interface
func (l *Light) TurnOn() error {
	return l.updateState("ON")
}

// TurnOff implements the switch interface
func (l *Light) TurnOff() error {
	return l.updateState("OFF")
}

// Toggle implements the Switch interface
func (l *Light) Toggle() error {
	if l.IsOn() {
		return l.TurnOff()
	}

	return l.TurnOn()
}

// Type implements the Component interface
func (l *Light) Type() components.Type {
	return components.TypeEsphomeLight
}

// Update implements the Component interface
func (l *Light) Update(value []byte) error {
	data := payload{}
	if err := json.Unmarshal(value, &data); err != nil {
		return err
	}

	if data.State != "" {
		l.On.Store(data.State == "ON")
	} else {
		l.On.Store(false)
	}

	l.Brightness.Store(uint32(data.Brightness))
	l.WarmWhite.Store(uint32(data.Color.W))
	l.ColdWhite.Store(uint32(data.Color.C))
	l.Red.Store(uint32(data.Color.R))
	l.Green.Store(uint32(data.Color.G))
	l.Blue.Store(uint32(data.Color.B))

	return nil
}

// Collectors implements the Component interface
func (l *Light) Collectors(labels prometheus.Labels) []prometheus.Collector {
	return []prometheus.Collector{
		components.GaugeCollector("esphome_light", labels,
			func() float64 {
				if l.IsOn() {
					return 1
				}
				return 0
			},
		),
	}
}
