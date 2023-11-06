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
		logger := controller.LoggerWithFields(log)

		current, err := controller.Temperature()
		if err != nil {
			logger.Warn("failed to get temperature", zap.Error(err))
			continue
		}

		if current == 0 {
			logger.Info("the temperature is reported to be 0, this is unlikely, let's ignore this for now")
			continue
		}

		target, err := controller.TemperatureTarget()
		if err != nil {
			logger.Warn("failed to get temperature target", zap.Error(err))
			continue
		}

		if target == 0 {
			logger.Info("the temperature target is 0, this is unlikely, let's ignore this for now")
			continue
		}

		min := target - b.Params.Hysteresis
		max := target + b.Params.Hysteresis

		logger = logger.With(
			zap.Float64("temperature", current),
			zap.Float64("target", target),
		)

		if controller.IsHeating() {
			if current > max {
				logger.Info("boiler entering the falling phase")
				err = controller.SetHeating(false)
			}
		} else {
			if current < min {
				logger.Info("boiler entering the rising phase")
				err = controller.SetHeating(true)
			}
		}
		if err != nil {
			logger.Warn("failed to set the component to heating", zap.Error(err))
		}

		shouldTurnOn = shouldTurnOn || controller.IsHeating()
	}

	logger := log.With(zap.Bool("state", shouldTurnOn))

	if b.IsOn() == shouldTurnOn {
		logger.Debug("boiler already in the good state")
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
