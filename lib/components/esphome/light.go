package esphome

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/gregdel/homed/lib/components"
	"github.com/gregdel/homed/lib/components/common"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	components.Register(components.TypeEsphomeLight, NewLight)
}

// Light represents a esphome light.
type Light struct {
	common.Component

	mu sync.RWMutex

	on         bool
	brightness uint32
	colorMode  string

	coldWhite uint32
	warmWhite uint32
	red       uint32
	green     uint32
	blue      uint32
}

type Snapshot struct {
	ID           string             `json:"id"`
	UpdatedAt    *time.Time         `json:"updated_at"`
	FriendlyName string             `json:"friendly_name"`
	Hide         bool               `json:"hide"`
	Device       *components.Device `json:"device"`
	On           bool               `json:"on"`
	Brightness   uint32             `json:"brightness"`
	ColorMode    string             `json:"color_mode"`
	ColdWhite    uint32             `json:"cold_white"`
	WarmWhite    uint32             `json:"warm_white"`
	Red          uint32             `json:"red"`
	Green        uint32             `json:"green"`
	Blue         uint32             `json:"blue"`
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
	l.mu.RLock()
	data := payload{
		State:      s,
		Brightness: uint8(l.brightness),
		ColorMode:  l.colorMode,
		Color: color{
			C: uint8(l.coldWhite),
			W: uint8(l.warmWhite),
		},
	}
	l.mu.RUnlock()

	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return l.WriteCommand(payload)
}

// IsOn implements the BinarySensor interface
func (l *Light) IsOn() bool {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return l.on
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

	l.mu.Lock()
	defer l.mu.Unlock()

	if data.State != "" {
		l.on = data.State == "ON"
	} else {
		l.on = false
	}

	l.brightness = uint32(data.Brightness)
	l.warmWhite = uint32(data.Color.W)
	l.coldWhite = uint32(data.Color.C)
	l.red = uint32(data.Color.R)
	l.green = uint32(data.Color.G)
	l.blue = uint32(data.Color.B)

	return nil
}

func (l *Light) Snapshot() any {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return Snapshot{
		ID:           l.ID(),
		UpdatedAt:    l.UpdatedAt.Load(),
		FriendlyName: l.FriendlyName(),
		Hide:         l.Hide,
		Device:       l.Device(),
		On:           l.on,
		Brightness:   l.brightness,
		ColorMode:    l.colorMode,
		ColdWhite:    l.coldWhite,
		WarmWhite:    l.warmWhite,
		Red:          l.red,
		Green:        l.green,
		Blue:         l.blue,
	}
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
