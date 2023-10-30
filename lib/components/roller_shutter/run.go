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
	dawn, dusk := rs.getDawnDusk(now)
	maxDelay := time.Duration(rs.Params.RandomDelay) * time.Minute
	return (now.Before(dawn) || now.After(dusk.Add(maxDelay)))
}

// nextEvent returns the time of the next event and the next state (open/close)
// as a bool.
func (rs *RollerShutter) nextEvent(shouldOpen bool) time.Time {
	now := time.Now()

	// Random minutes between 0 and random delay
	s := rand.NewSource(now.Unix())
	rd := rand.New(s)

	var randomDelay time.Duration
	if rs.Params.RandomDelay != 0 {
		randomDelay = time.Duration(
			rd.Intn(rs.Params.RandomDelay)) * time.Minute
	}
	maxDelay := time.Duration(rs.Params.RandomDelay) * time.Minute

	dawn, dusk := rs.getDawnDusk(now)
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
		rs.logger.Error("should never reach this code")
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
	nextDawn, _ := rs.getDawnDusk(now.Add(24 * time.Hour))
	return nextDawn.Add(-1 * randomDelay)
}
