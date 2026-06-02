package esphome

import (
	"encoding/json"
	"sync"

	"github.com/gregdel/homed/lib/components"
	"github.com/gregdel/homed/lib/components/common"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	components.Register(components.TypeEsphomeLight, NewLight)
}

const (
	defaultColdWhite = 128
	defaultWarmWhite = 127
)

// Light represents a esphome light.
type Light struct {
	common.Component

	mu sync.RWMutex

	on         bool
	brightness uint32
	colorMode  string
	colorTemp  uint32

	coldWhite uint32
	warmWhite uint32
	red       uint32
	green     uint32
	blue      uint32
}

type Snapshot struct {
	common.SnapshotBase
	On         bool   `json:"on"`
	Brightness uint32 `json:"brightness"`
	ColorMode  string `json:"color_mode"`
	ColorTemp  uint32 `json:"color_temp"`
	ColdWhite  uint32 `json:"cold_white"`
	WarmWhite  uint32 `json:"warm_white"`
	Red        uint32 `json:"red"`
	Green      uint32 `json:"green"`
	Blue       uint32 `json:"blue"`
}

// NewLight returns a new light
func NewLight() components.Component {
	return &Light{}
}

type color struct {
	C *uint8 `json:"c,omitempty"`
	W *uint8 `json:"w,omitempty"`
	R *uint8 `json:"r,omitempty"`
	G *uint8 `json:"g,omitempty"`
	B *uint8 `json:"b,omitempty"`
}

type payload struct {
	State      string  `json:"state,omitempty"`
	Brightness *uint8  `json:"brightness,omitempty"`
	ColorMode  string  `json:"color_mode,omitempty"`
	ColorTemp  *uint16 `json:"color_temp,omitempty"`
	Color      *color  `json:"color,omitempty"`
}

func uint8Ptr(value uint8) *uint8 {
	return &value
}

// TurnOn implements the switch interface
func (l *Light) updateState(s string) error {
	l.mu.RLock()
	brightness := uint8(l.brightness)
	colorMode := l.colorMode
	coldWhite := uint8(l.coldWhite)
	warmWhite := uint8(l.warmWhite)
	if colorMode != "rgb" && colorMode != "cwww" {
		colorMode = "cwww"
	}
	if coldWhite == 0 && warmWhite == 0 {
		coldWhite = defaultColdWhite
		warmWhite = defaultWarmWhite
	}

	data := payload{
		State:      s,
		Brightness: &brightness,
		ColorMode:  colorMode,
	}
	switch colorMode {
	case "rgb":
		data.Color = &color{
			R: uint8Ptr(uint8(l.red)),
			G: uint8Ptr(uint8(l.green)),
			B: uint8Ptr(uint8(l.blue)),
		}
	default:
		data.Color = &color{
			C: uint8Ptr(coldWhite),
			W: uint8Ptr(warmWhite),
		}
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
	}

	if data.Brightness != nil {
		l.brightness = uint32(*data.Brightness)
	}
	if data.ColorMode != "" {
		l.colorMode = data.ColorMode
	}
	if data.ColorTemp != nil {
		l.colorTemp = uint32(*data.ColorTemp)
	}
	if data.Color != nil {
		if data.Color.W != nil {
			l.warmWhite = uint32(*data.Color.W)
		}
		if data.Color.C != nil {
			l.coldWhite = uint32(*data.Color.C)
		}
		if data.Color.R != nil {
			l.red = uint32(*data.Color.R)
		}
		if data.Color.G != nil {
			l.green = uint32(*data.Color.G)
		}
		if data.Color.B != nil {
			l.blue = uint32(*data.Color.B)
		}
	}

	return nil
}

func (l *Light) ValuesSnapshot() any {
	l.mu.RLock()
	on := l.on
	brightness := l.brightness
	colorMode := l.colorMode
	colorTemp := l.colorTemp
	coldWhite := l.coldWhite
	warmWhite := l.warmWhite
	red := l.red
	green := l.green
	blue := l.blue
	l.mu.RUnlock()

	return Snapshot{
		SnapshotBase: l.SnapshotBase(),
		On:           on,
		Brightness:   brightness,
		ColorMode:    colorMode,
		ColorTemp:    colorTemp,
		ColdWhite:    coldWhite,
		WarmWhite:    warmWhite,
		Red:          red,
		Green:        green,
		Blue:         blue,
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
