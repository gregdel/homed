package boiler

import (
	"context"
	"log/slog"

	"github.com/gregdel/homed/lib/components"
)

// Run implements the Component interface
func (b *Boiler) Run(ctx context.Context, logger *slog.Logger, inventory *components.Components) error {
	b.Events.Incoming = make(chan components.Event)

	log := logger.With(slog.String("device", "boiler"))

	if err := b.YAMLParams.Decode(&b.Params); err != nil {
		return err
	}

	for _, component := range inventory.List() {
		if component.Type() != components.TypeHomedTemperature {
			continue
		}

		tempInternal := component.(components.TemperatureControllerInternal)
		b.controllers = append(b.controllers, tempInternal)

		tempInternal.Subscribe(b.ID(), b.Events.Incoming)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case e := <-b.Events.Incoming:
			b.checkState(log.With(slog.String("event_id", e.ID)))
		}
	}
}

func (b *Boiler) checkState(log *slog.Logger) {
	shouldTurnOn := false
	for _, controller := range b.controllers {
		shouldTurnOn = shouldTurnOn || controller.IsHeating()
	}

	logger := log.With(slog.Bool("state", shouldTurnOn))

	if b.IsOn() == shouldTurnOn {
		return
	}

	logger.Info("changing state")

	var err error
	if shouldTurnOn {
		err = b.TurnOn()
	} else {
		err = b.TurnOff()
	}
	if err != nil {
		logger.Warn("failed to set state", slog.Any("error", err))
	}
}
