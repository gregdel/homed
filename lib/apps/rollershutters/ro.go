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

const (
	name        = "roller_shutters"
	randomDelay = 20
)

func init() {
	apps.Register(app())
}

type rollerShutters struct {
	rs components.RollerShutter

	config *config.Config
}

func app() *rollerShutters {
	return &rollerShutters{}
}

func (r *rollerShutters) Name() string {
	return name
}

func (r *rollerShutters) Init(config *config.Config) error {
	r.config = config
	return nil
}

func (r *rollerShutters) getSunriseSunset(config *apps.Config, t time.Time) (time.Time, time.Time) {
	_, utcOffset := t.Zone()
	p := sunrisesunset.Parameters{
		Latitude:  r.config.Location.Latitude,
		Longitude: r.config.Location.Longitude,
		UtcOffset: float64(utcOffset / 3600),
		Date:      t,
	}

	sunrise, sunset, err := p.GetSunriseSunset()
	if err != nil {
		config.Logger.Error(
			"failed to get sunrise / sunset times, falling back",
			zap.Error(err),
		)

		// Fallback to 8:00 (sunrise) / 19:00 (sunset)
		sunrise = time.Date(t.Year(), t.Month(), t.Day(), 8, 00, 0, 0, t.Location())
		sunset = time.Date(t.Year(), t.Month(), t.Day(), 19, 0, 0, 0, t.Location())
	}

	return sunrise, sunset
}

// nextEvent returns the time of the next event and the next state (open/close)
// as a bool.
func (r *rollerShutters) nextEvent(config *apps.Config) (time.Time, bool) {
	now := time.Now()

	// Random minutes between 0 and random delay
	s := rand.NewSource(now.Unix())
	rd := rand.New(s)
	minutes := time.Duration(rd.Intn(randomDelay)) * time.Minute

	sunrise, sunset := r.getSunriseSunset(config, now)

	// Open the roller shutter 0 to X minutes after the sunrise
	openTime := sunrise.Add(-1 * minutes)
	if now.Before(openTime) {
		config.Logger.Info(
			"settting roller shutter open time",
			zap.Time("sunrise", sunrise),
			zap.Duration("offset", -1*minutes),
		)
		return openTime, true
	}

	// Close the roller shutter X minutes after the sunset
	closeTime := sunset.Add(minutes)
	if now.Before(closeTime) {
		config.Logger.Info(
			"settting roller close time",
			zap.Time("sunset", sunset),
			zap.Duration("offset", minutes),
		)
		return closeTime, false
	}

	// We're after the sunset, get the sunrise of the next morning
	sunrise, _ = r.getSunriseSunset(config, now.Add(24*time.Hour))
	config.Logger.Info(
		"settting roller shutter open time to the next day",
		zap.Time("sunrise", sunrise),
		zap.Duration("offset", -1*minutes),
	)
	return sunrise.Add(-1 * minutes), true
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
	if r.rs == nil {
		if err := r.update(config); err != nil {
			return err
		}
	}

	for {
		event, open := r.nextEvent(config)
		duration := event.Sub(time.Now())

		config.Logger.Info("setting next event",
			zap.Time("next_event", event),
			zap.Duration("sleep_duration", duration),
		)

		select {
		case <-ctx.Done():
			return nil
		case <-time.After(duration):
			if r.rs.IsOpen() == open {
				config.Logger.Info("roller shutter already in good state")
				continue
			}

			var err error
			if open {
				config.Logger.Info("openning the roller shutters")
				err = r.rs.Open()
			} else {
				config.Logger.Info("closing the roller shutter")
				err = r.rs.Close()
			}

			if err != nil {
				config.Logger.Info("failed to change the roller shutter state", zap.Error(err))
			}
		}
	}
}
