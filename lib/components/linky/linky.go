package linky

import (
	"encoding/json"

	"github.com/gregdel/homed/lib/components"
	"github.com/gregdel/homed/lib/components/common"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	components.Register(components.TypeLinky, New)
}

// Linky represents a ZLinky_TIC device over zigbee2mqtt
type Linky struct {
	common.Component

	CurrentPower int `json:"apparent_power"`
	CounterHC    int `json:"current_tier1_summ_delivered"`
	CounterHP    int `json:"current_tier2_summ_delivered"`
}

// New returns a new component
func New() components.Component {
	return &Linky{}
}

// Type implements the Component interface
func (l *Linky) Type() components.Type {
	return components.TypeLinky
}

// Collectors implements the Component interface
func (l *Linky) Collectors(labels prometheus.Labels) []prometheus.Collector {
	return []prometheus.Collector{
		components.CounterCollector("linky_hp", labels,
			func() float64 { return float64(l.CounterHP) },
		),
		components.CounterCollector("linky_hc", labels,
			func() float64 { return float64(l.CounterHC) },
		),
		components.GaugeCollector("power", labels,
			func() float64 { return float64(l.CurrentPower) },
		),
	}
}

// Update implements the components interface
func (l *Linky) Update(data []byte) error {
	return json.Unmarshal(data, &l)
}
