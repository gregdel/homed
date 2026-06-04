package rollershutter

import (
	"log/slog"
	"time"

	"github.com/gregdel/homed/lib/components"
	"github.com/gregdel/homed/lib/components/common"
	"github.com/gregdel/homed/lib/config"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	components.Register(components.TypeRollerShutter, NewRollerShutter)
}

// Params represents the roller shutter params
type Params struct {
	Enabled     bool            `yaml:"enabled"`
	Location    config.Location `yaml:"location"`
	OpenAfter   string          `yaml:"open_after"`
	OpenBefore  string          `yaml:"open_before"`
	CloseAfter  string          `yaml:"close_after"`
	CloseBefore string          `yaml:"close_before"`
}

type Snapshot struct {
	common.SnapshotBase
	Value     float64            `json:"value"`
	NextEvent *NextEventSnapshot `json:"next_event"`
}

type NextEventSnapshot struct {
	Action      string    `json:"action"`
	ScheduledAt time.Time `json:"scheduled_at"`
}

// RollerShutter represents a generic roller shutter
type RollerShutter struct {
	common.GenericSensor
	Params      Params
	openWindow  dailyWindow
	closeWindow dailyWindow
	logger      *slog.Logger
}

// NewRollerShutter returns a new cover component
func NewRollerShutter() components.Component {
	return &RollerShutter{}
}

// Type implements the Component interface
func (rs *RollerShutter) Type() components.Type {
	return components.TypeRollerShutter
}

func (rs *RollerShutter) ValidateConfig() error {
	return rs.setup()
}

func (rs *RollerShutter) ValuesSnapshot() any {
	return rs.valuesSnapshotAt(time.Now())
}

func (rs *RollerShutter) valuesSnapshotAt(now time.Time) Snapshot {
	return Snapshot{
		SnapshotBase: rs.SnapshotBase(),
		Value:        rs.SensorValue(),
		NextEvent:    rs.nextEventSnapshotAt(now),
	}
}

func (rs *RollerShutter) nextEventSnapshotAt(now time.Time) *NextEventSnapshot {
	if !rs.Params.Enabled {
		return nil
	}

	nextEvent := rs.nextEvent(now)
	return &NextEventSnapshot{
		Action:      nextEvent.action.String(),
		ScheduledAt: nextEvent.scheduledAt,
	}
}

// Collectors implements the Component interface
func (rs *RollerShutter) Collectors(labels prometheus.Labels) []prometheus.Collector {
	return nil
}

// OpenedAt implements the RollerShutter interface
func (rs *RollerShutter) OpenedAt() float64 {
	return rs.SensorValue()
}

// IsOpen implements the RollerShutter interface
func (rs *RollerShutter) IsOpen() bool {
	return rs.SensorValue() == 100
}

// IsClosed implements the RollerShutter interface
func (rs *RollerShutter) IsClosed() bool {
	return !rs.IsOpen()
}

// Open implements the RollerShutter interface
func (rs *RollerShutter) Open() error {
	return rs.WriteCommand([]byte("open"))
}

// Close implements the RollerShutter interface
func (rs *RollerShutter) Close() error {
	return rs.WriteCommand([]byte("close"))
}

// Stop implements the RollerShutter interface
func (rs *RollerShutter) Stop() error {
	return rs.WriteCommand([]byte("stop"))
}
