package homedtemperature

import (
	"context"

	"github.com/gregdel/homed/lib/components"
	"github.com/gregdel/homed/lib/components/common"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

func init() {
	components.Register(components.TypeHomedTemperatureSwitch, NewTemperatureSwitch)
}

// NewTemperatureSwitch returns a new general switch to turn on and off the
// homed temperature components
func NewTemperatureSwitch() components.Component {
	return &TemperatureSwitch{}
}

// TemperatureSwitch represents a virtual button that can trigger multiple switches
type TemperatureSwitch struct {
	common.VirtualSwitch
}

// Type implements the Component interface
func (sw *TemperatureSwitch) Type() components.Type {
	return components.TypeHomedTemperatureSwitch
}

// Run implements the Component interface
func (sw *TemperatureSwitch) Run(ctx context.Context, logger *zap.Logger,
	inventory *components.Components) error {
	sw.Device().Online.Store(true)
	return sw.VirtualSwitch.Run(ctx, logger, inventory)
}

// Collectors implements the Component interface
func (sw *TemperatureSwitch) Collectors(labels prometheus.Labels) []prometheus.Collector {
	return []prometheus.Collector{
		components.CounterCollector("temperature_switch", labels,
			func() float64 { return sw.Counter.Load() },
		),
	}
}
