package ro

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gregdel/homed/lib/apps"
	"github.com/gregdel/homed/lib/components"
	"github.com/gregdel/homed/lib/config"
	"go.uber.org/zap"
)

const name = "roller_shutters"

func init() {
	apps.Register(app())
}

type rollerShutters struct {
	rs components.RollerShutter
}

func app() *rollerShutters {
	return &rollerShutters{}
}

func (r *rollerShutters) Name() string {
	return name
}

func (r *rollerShutters) Init(config *config.Config) error {
	return nil
}

func (r *rollerShutters) nextEvent() (time.Time, bool) {
	now := time.Now()
	s := rand.NewSource(now.Unix())
	rd := rand.New(s)

	// Random minutes between 0 and 30
	minutes := time.Duration(rd.Intn(30)) * time.Minute

	// sunrise (open)
	openStart := time.Date(
		now.Year(), now.Month(), now.Day(),
		8, 10, 0, 0,
		now.Location(),
	)

	// We are before the start of the opening event
	if now.Before(openStart) {
		return openStart.Add(minutes), true
	}

	// sunset (close)
	closeStart := time.Date(
		now.Year(), now.Month(), now.Day(),
		16, 55, 0, 0,
		now.Location(),
	)

	if now.Before(closeStart) {
		return closeStart.Add(minutes), false
	}

	return openStart.Add(minutes).Add(24 * time.Hour), true
}

func (r *rollerShutters) Run(ctx *apps.RunCtx) error {
	for _, component := range ctx.Components.List() {
		if component.Type() != components.TypeRollerShutter {
			continue
		}

		r.rs = component.(components.RollerShutter)
		break
	}

	if r.rs == nil {
		return fmt.Errorf("roller_shutters: failed to find component")
	}

	for {
		event, open := r.nextEvent()
		duration := event.Sub(time.Now())

		ctx.Logger.Info("setting next event",
			zap.Time("next_event", event),
			zap.Duration("sleep_duration", duration),
		)

		select {
		case <-ctx.Ctx.Done():
			return nil
		case <-time.After(duration):
			if r.rs.IsOpen() == open {
				ctx.Logger.Info("roller shutter already in good state")
				continue
			}

			var err error
			if open {
				ctx.Logger.Info("openning the roller shutters")
				err = r.rs.Open()
			} else {
				ctx.Logger.Info("closing the roller shutter")
				err = r.rs.Close()
			}

			if err != nil {
				ctx.Logger.Info("failed to change the roller shutter state", zap.Error(err))
			}
		}
	}
}
