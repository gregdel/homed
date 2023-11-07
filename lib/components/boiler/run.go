package boiler

import (
	"context"

	"github.com/gregdel/homed/lib/components"
	"go.uber.org/zap"
)

// Run implements the Component interface
func (b *Boiler) Run(ctx context.Context, logger *zap.Logger, inventory *components.Components) error {
	b.Events.Incoming = make(chan components.Event)

	log := b.LoggerWithFields(logger)

	if err := b.YAMLParams.Decode(&b.Params); err != nil {
		return err
	}

	for _, component := range inventory.List() {
		if component.Type() != components.TypeHomedTemperature {
			continue
		}

		tempInternal := component.(components.TemperatureControllerInternal)
		b.controllers = append(b.controllers, tempInternal)
		tempInternal.SetHeating(false)

		tempInternal.Subscribe(b.ID(), b.Events.Incoming)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-b.Events.Incoming:
			b.checkState(log)
		}
	}
}

func (b *Boiler) checkState(log *zap.Logger) {
	shouldTurnOn := false
	for _, controller := range b.controllers {
		shouldTurnOn = shouldTurnOn || controller.IsHeating()
	}

	logger := log.With(zap.Bool("state", shouldTurnOn))

	if b.IsOn() == shouldTurnOn {
		return
	}

	logger.Info("changing boiler state")

	var err error
	if shouldTurnOn {
		err = b.TurnOn()
	} else {
		err = b.TurnOff()
	}
	if err != nil {
		logger.Warn("failed to set boiler state", zap.Error(err))
	}
}
