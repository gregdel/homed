package rollershutter

import (
	"context"
	"math/rand"
	"time"

	"go.uber.org/zap"
)

// Run implements the component interface
func (rs *RollerShutter) Run(ctx context.Context, logger *zap.Logger) error {
	if err := rs.YAMLParams.Decode(&rs.Params); err != nil {
		return err
	}

	log := logger.With(zap.String("component_id", rs.ID()))
	rs.logger = log

	if !rs.Params.Enabled {
		log.Info("automatic control disabled")
		return nil
	}

	shouldOpen := rs.isNextEventOpen()

	for {
		event := rs.nextEvent(shouldOpen)
		duration := time.Until(event)

		log.Info("setting next event",
			zap.Time("next_event", event),
			zap.Bool("should_open", shouldOpen),
			zap.Duration("sleep_duration", duration),
		)

		select {
		case <-ctx.Done():
			log.Info("stopping")
			return nil
		case <-time.After(duration):
			var err error
			if shouldOpen {
				log.Info("openning the roller shutters")
				err = rs.Open()
			} else {
				log.Info("closing the roller shutter")
				err = rs.Close()
			}

			shouldOpen = !shouldOpen

			if err != nil {
				log.Info("failed to change the roller shutter state", zap.Error(err))
			}
		}
	}
}

func (rs *RollerShutter) isNextEventOpen() bool {
	now := time.Now()
	sunrise, sunset := rs.getSunriseSunset(now)
	maxDelay := time.Duration(rs.Params.RandomDelay) * time.Minute
	return (now.Before(sunrise) || now.After(sunset.Add(maxDelay)))
}

// nextEvent returns the time of the next event and the next state (open/close)
// as a bool.
func (rs *RollerShutter) nextEvent(shouldOpen bool) time.Time {
	now := time.Now()
	sunrise, sunset := rs.getSunriseSunset(now)

	// Random minutes between 0 and random delay
	s := rand.NewSource(now.Unix())
	rd := rand.New(s)

	var randomDelay, maxDelay time.Duration
	if rs.Params.RandomDelay != 0 {
		randomDelay = time.Duration(
			rd.Intn(rs.Params.RandomDelay),
		) * time.Minute
		maxDelay = time.Duration(
			rs.Params.RandomDelay,
		) * time.Minute
	}

	if !shouldOpen {
		if now.Before(sunset) {
			return sunset.Add(randomDelay)
		}

		// In the gray area between sunset and (sunset + max delay), return
		// (sunset + max delay)
		if now.Before(sunset.Add(maxDelay)) {
			return sunset.Add(maxDelay)
		}

		// We should never reach this point
		rs.logger.Error("should never reach this code")
	}

	openStart := sunrise.Add(-1 * maxDelay)
	openEnd := sunrise

	// Before the (sunrise + max delay)
	if now.Before(openStart) {
		return sunrise.Add(-1 * randomDelay)
	}

	// In the gray area between the (sunrise - max delay) and sunrise, make
	// sure we always return the sunrise time
	if now.After(openStart) && now.Before(openEnd) {
		return sunrise
	}

	// We're after sunset, get the sunrise of the next morning
	sunriseNextDay, _ := rs.getSunriseSunset(now.Add(24 * time.Hour))
	return sunriseNextDay.Add(-1 * randomDelay)
}
