package ro

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/gregdel/homed/lib/apps"
	"github.com/gregdel/homed/lib/components"
	"github.com/gregdel/homed/lib/config"
	"github.com/kelvins/sunrisesunset"
	"go.uber.org/zap"
)

func init() {
	apps.Register(app())
}

type rollerShutters struct {
	// TODO: handle multiple roller shutters for @PouuleT
	rs components.RollerShutter

	logger *zap.Logger
	config *config.Config
}

func app() *rollerShutters {
	return &rollerShutters{}
}

func (r *rollerShutters) Name() string {
	return "roller_shutters"
}

func (r *rollerShutters) Init(config *config.Config) error {
	r.config = config
	return nil
}

func (r *rollerShutters) getDawnDusk(t time.Time) (time.Time, time.Time) {
	sunrise, sunset := r.getSunriseSunset(t)

	// Let's say the difference between dawn -> sunrise and sunset -> dusk is
	// around 40 minutes.
	d := 40 * time.Minute

	return sunrise.Add(-1 * d), sunset.Add(d)
}

func (r *rollerShutters) getSunriseSunset(t time.Time) (time.Time, time.Time) {
	_, utcOffset := t.Zone()
	p := sunrisesunset.Parameters{
		Latitude:  r.config.Location.Latitude,
		Longitude: r.config.Location.Longitude,
		UtcOffset: float64(utcOffset / 3600),
		Date:      t,
	}

	sunrise, sunset, err := p.GetSunriseSunset()
	if err != nil {
		r.logger.Error(
			"failed to get sunrise / sunset times, falling back",
			zap.Error(err),
		)

		// Fallback to 8:00 (sunrise) / 19:00 (sunset)
		sunrise = time.Date(t.Year(), t.Month(), t.Day(), 8, 00, 0, 0, t.Location())
		sunset = time.Date(t.Year(), t.Month(), t.Day(), 19, 0, 0, 0, t.Location())
	}

	return sunrise, sunset
}

func (r *rollerShutters) isNextEventOpen() bool {
	now := time.Now()
	dawn, dusk := r.getDawnDusk(now)
	maxDelay := time.Duration(r.config.RollerShutter.RandomDelay) * time.Minute
	return (now.Before(dawn) || now.After(dusk.Add(maxDelay)))
}

// nextEvent returns the time of the next event and the next state (open/close)
// as a bool.
func (r *rollerShutters) nextEvent(shouldOpen bool) time.Time {
	now := time.Now()

	// Random minutes between 0 and random delay
	s := rand.NewSource(now.Unix())
	rd := rand.New(s)

	var randomDelay time.Duration
	if r.config.RollerShutter.RandomDelay != 0 {
		randomDelay = time.Duration(
			rd.Intn(r.config.RollerShutter.RandomDelay)) * time.Minute
	}
	maxDelay := time.Duration(r.config.RollerShutter.RandomDelay) * time.Minute

	dawn, dusk := r.getDawnDusk(now)
	if !shouldOpen {
		// Should close between dusk and dusk + random delay
		if now.Before(dusk) {
			return dusk.Add(randomDelay)
		}

		// In the gray area between dusk and dusk + max delay, return dusk + max delay
		if now.Before(dusk.Add(maxDelay)) {
			return dusk.Add(maxDelay)
		}

		// We should never reach this point
		r.logger.Error("should never reach this code")
	}

	openStart := dawn.Add(-1 * maxDelay)
	openEnd := dawn

	// Before the dawn + random delay
	if now.Before(openStart) {
		return dawn.Add(-1 * randomDelay)
	}

	// In the gray area between the dawn - random delay and dawn, make sure
	// we always return the dawn time
	if now.After(openStart) && now.Before(openEnd) {
		// Make sure we don't add a random delay in this range
		return dawn
	}

	// We're after dusk, get the dawn of the next morning
	nextDawn, _ := r.getDawnDusk(now.Add(24 * time.Hour))
	return nextDawn.Add(-1 * randomDelay)
}

// update finds the roller shutter in the component list
func (r *rollerShutters) update(config *apps.Config) error {
	for _, component := range config.Components.List() {
		if component.Type() != components.TypeRollerShutter {
			continue
		}

		r.rs = component.(components.RollerShutter)
		return nil
	}

	return fmt.Errorf("roller_shutters: failed to find component")
}

func (r *rollerShutters) Run(ctx context.Context, config *apps.Config) error {
	logger := config.Logger.With(zap.String("app", r.Name()))

	if !r.config.RollerShutter.Enabled {
		logger.Info("app is disabled")
		return nil
	}

	r.logger = logger

	if r.rs == nil {
		if err := r.update(config); err != nil {
			return err
		}
	}

	shouldOpen := r.isNextEventOpen()

	for {
		event := r.nextEvent(shouldOpen)
		duration := time.Until(event)

		logger.Info("setting next event",
			zap.Time("next_event", event),
			zap.Bool("should_open", shouldOpen),
			zap.Duration("sleep_duration", duration),
		)

		select {
		case <-ctx.Done():
			return nil
		case <-time.After(duration):
			var err error
			if shouldOpen {
				logger.Info("openning the roller shutters")
				err = r.rs.Open()
			} else {
				logger.Info("closing the roller shutter")
				err = r.rs.Close()
			}

			shouldOpen = !shouldOpen

			if err != nil {
				logger.Info("failed to change the roller shutter state", zap.Error(err))
			}
		}
	}
}
