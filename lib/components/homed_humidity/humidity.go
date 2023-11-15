package homedhumidity

import (
	"context"
	"fmt"

	"github.com/gregdel/homed/lib/components"
	"github.com/gregdel/homed/lib/components/common"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

func init() {
	components.Register(components.TypeHomedHumidity, New)
}

// Params are the params of the humidity controller
type Params struct {
	Sensor            string  `yaml:"sensor"`
	Switch            string  `yaml:"switch"`
	HumidityThreshold float64 `yaml:"humidity_threshold"`
}

// HomedHumidity is a component that controls a switch based on a humidity
// sensor.
type HomedHumidity struct {
	common.Component
	Params Params
}

// New returns a new homed humidity component
func New() components.Component {
	return &HomedHumidity{}
}

// Type implements the Component interface
func (h *HomedHumidity) Type() components.Type {
	return components.TypeHomedHumidity
}

// Collectors implements the Component interface
func (h *HomedHumidity) Collectors(labels prometheus.Labels) []prometheus.Collector {
	return nil
}

// Update implements the Component interface
func (h *HomedHumidity) Update([]byte) error {
	return nil
}

// Run implements the Component interface
func (h *HomedHumidity) Run(ctx context.Context, logger *zap.Logger, inventory *components.Components) error {
	if err := h.YAMLParams.Decode(&h.Params); err != nil {
		return err
	}

	log := logger.With(
		zap.String("type", "humidity_controller"),
		zap.String("room", h.Device().Room),
	)

	sensorComponent, err := inventory.Get(h.Params.Sensor)
	if err != nil {
		log.Error("failed to get sensor", zap.Error(err))
	}

	sensor, ok := (sensorComponent).(components.HumidityGetter)
	if !ok {
		return fmt.Errorf("component %q is not a humidity getter", h.Params.Sensor)
	}

	swComponent, err := inventory.Get(h.Params.Switch)
	if err != nil {
		log.Error("failed to get switch", zap.Error(err))
	}

	sw, ok := (swComponent).(components.Switch)
	if !ok {
		return fmt.Errorf("component %q is not a switch", h.Params.Switch)
	}

	h.Events.Incoming = make(chan components.Event)
	sensor.Subscribe(h.ID(), h.Events.Incoming)

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-h.Events.Incoming:
			humidity, err := sensor.Humidity()
			if err != nil {
				log.Warn("failed to get humidity", zap.Error(err))
				continue
			}

			if humidity > h.Params.HumidityThreshold {
				log.Info("FAN should be ON")
				sw.TurnOn()
			} else {
				log.Info("FAN should be OFF")
				sw.TurnOff()
			}
		}
	}
}
