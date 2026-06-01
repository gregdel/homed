package boiler

import (
	"errors"
	"sync"
	"time"

	"github.com/gregdel/homed/lib/components"
	"github.com/gregdel/homed/lib/components/common"
	"github.com/gregdel/homed/lib/config"
	"github.com/prometheus/client_golang/prometheus"
)

const cooldownDuration = 5 * time.Minute

var ErrBoilerCooldown = errors.New("components: boiler: last change is too recent")

func init() {
	components.Register("boiler", NewBoiler)
}

type Snapshot struct {
	common.SnapshotBase
	On              bool       `json:"on"`
	LastStateChange *time.Time `json:"last_state_change"`
}

// Boiler is a component that controls the boiler
type Boiler struct {
	common.Switch

	mu sync.RWMutex

	Params      config.TemperatureControl
	controllers []components.TemperatureControllerInternal

	lastStateChange time.Time
}

// NewBoiler returns a new status component
func NewBoiler() components.Component {
	return &Boiler{
		controllers: []components.TemperatureControllerInternal{},
	}
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
				if b.IsOn() {
					return 1
				}
				return 0
			},
		),
	}
}

func (b *Boiler) ValuesSnapshot() any {
	return Snapshot{
		SnapshotBase:    b.SnapshotBase(),
		On:              b.IsOn(),
		LastStateChange: b.lastStateChangePtr(),
	}
}

func (b *Boiler) lastStateChangePtr() *time.Time {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.lastStateChange.IsZero() {
		return nil
	}

	t := b.lastStateChange
	return &t
}

func (b *Boiler) reserveStateChange(now time.Time) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if !b.lastStateChange.IsZero() && now.Before(b.lastStateChange.Add(cooldownDuration)) {
		return ErrBoilerCooldown
	}

	b.lastStateChange = now
	return nil
}

// WriteCommand implements the Component interface
func (b *Boiler) WriteCommand(data []byte) error {
	if err := b.reserveStateChange(time.Now()); err != nil {
		return err
	}

	return b.Switch.WriteCommand(data)
}
