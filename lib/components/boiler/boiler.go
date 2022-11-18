package boiler

import (
	"fmt"
	"time"

	"github.com/gregdel/homed/lib/components"
	"github.com/gregdel/homed/lib/components/common"
	"github.com/prometheus/client_golang/prometheus"
)

const cooldownDuration = 5 * time.Minute

func init() {
	components.Register("boiler", NewBoiler)
}

// Boiler is a component that controls the boiler
type Boiler struct {
	common.Switch

	LastStateChange *time.Time `json:"last_state_change"`
}

// NewBoiler returns a new status component
func NewBoiler() components.Component {
	return &Boiler{}
}

// Type implements the Component interface
func (b *Boiler) Type() components.Type {
	return components.TypeBoiler
}

// Collectors implements the Component interface
func (b *Boiler) Collectors(labels prometheus.Labels) []prometheus.Collector {
	return []prometheus.Collector{
		components.GaugeCollector("boiler", labels,
			func() float64 {
				if b.On {
					return 1
				}
				return 0
			},
		),
	}
}

// WriteCommand implements the Component interface
func (b *Boiler) WriteCommand(data []byte) error {
	now := time.Now()
	if b.LastStateChange == nil {
		b.LastStateChange = &now
	} else if b.LastStateChange.Add(cooldownDuration).Before(now) {
		return fmt.Errorf("components: boiler: last change is to recent")
	}

	return b.Switch.WriteCommand(data)
}
