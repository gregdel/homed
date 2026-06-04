package rollershutter

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/gregdel/homed/lib/components"
)

// Run implements the component interface
func (rs *RollerShutter) Run(ctx context.Context, logger *slog.Logger, _ *components.Components) error {
	log := logger.With(slog.String("component_id", rs.ID()))
	rs.logger = log

	if !rs.Params.Enabled {
		log.Info("automatic control disabled")
		return nil
	}

	for {
		nextEvent := rs.nextEvent(time.Now())
		sleepDuration := time.Until(nextEvent.scheduledAt)

		log.Info("setting next roller shutter action",
			slog.Time("run_at", nextEvent.scheduledAt),
			slog.String("action", nextEvent.action.String()),
			slog.Duration("sleep_duration", sleepDuration),
		)

		select {
		case <-ctx.Done():
			log.Info("stopping")
			return nil
		case <-time.After(sleepDuration):
			action := nextEvent.action
			err := rs.performAction(action)
			if err != nil {
				log.Error("failed to run roller shutter action",
					slog.String("action", action.String()),
					slog.Any("error", err),
				)
			}
			rs.Notify()
		}
	}
}

type scheduledEvent struct {
	action      shutterAction
	scheduledAt time.Time
}

type shutterAction int

const (
	openShutters shutterAction = iota
	closeShutters
)

func (action shutterAction) String() string {
	switch action {
	case openShutters:
		return "open"
	case closeShutters:
		return "close"
	default:
		return fmt.Sprintf("unknown(%d)", action)
	}
}

func (rs *RollerShutter) performAction(action shutterAction) error {
	switch action {
	case openShutters:
		rs.logger.Info("opening the roller shutters")
		return rs.Open()
	case closeShutters:
		rs.logger.Info("closing the roller shutters")
		return rs.Close()
	default:
		return fmt.Errorf("unknown roller shutter action: %s", action)
	}
}

func (rs *RollerShutter) nextEvent(now time.Time) scheduledEvent {
	openAt := rs.scheduledOpenTimeAt(now)
	if !now.After(openAt) {
		return scheduledEvent{action: openShutters, scheduledAt: openAt}
	}

	closeAt := rs.scheduledCloseTimeAt(now)
	if !now.After(closeAt) {
		return scheduledEvent{action: closeShutters, scheduledAt: closeAt}
	}

	nextOpenAt := rs.scheduledOpenTimeAt(now.AddDate(0, 0, 1))
	return scheduledEvent{action: openShutters, scheduledAt: nextOpenAt}
}

func (rs *RollerShutter) scheduledOpenTimeAt(now time.Time) time.Time {
	sunrise, _ := rs.getSunriseSunset(now)

	return clampToDailyWindow(sunrise.In(now.Location()), now, rs.openWindow)
}

func (rs *RollerShutter) scheduledCloseTimeAt(now time.Time) time.Time {
	_, sunset := rs.getSunriseSunset(now)

	return clampToDailyWindow(sunset.In(now.Location()), now, rs.closeWindow)
}

func (rs *RollerShutter) setup() error {
	if err := rs.YAMLParams.Decode(&rs.Params); err != nil {
		return err
	}

	return rs.parseDailyWindows()
}

func (rs *RollerShutter) parseDailyWindows() error {
	openWindow, err := parseDailyWindow("open", rs.Params.OpenAfter, rs.Params.OpenBefore)
	if err != nil {
		return err
	}

	closeWindow, err := parseDailyWindow("close", rs.Params.CloseAfter, rs.Params.CloseBefore)
	if err != nil {
		return err
	}

	rs.openWindow = openWindow
	rs.closeWindow = closeWindow
	return nil
}

type dailyWindow struct {
	start *clockTime
	end   *clockTime
}

type clockTime struct {
	hour   int
	minute int
}

func (ct clockTime) before(other clockTime) bool {
	if ct.hour != other.hour {
		return ct.hour < other.hour
	}

	return ct.minute < other.minute
}

func (ct clockTime) onDate(date time.Time) time.Time {
	return time.Date(
		date.Year(),
		date.Month(),
		date.Day(),
		ct.hour,
		ct.minute,
		0,
		0,
		date.Location(),
	)
}

func parseDailyWindow(name, afterValue, beforeValue string) (dailyWindow, error) {
	start, err := parseClockTime(afterValue)
	if err != nil {
		return dailyWindow{}, fmt.Errorf("invalid %s_after: %w", name, err)
	}

	end, err := parseClockTime(beforeValue)
	if err != nil {
		return dailyWindow{}, fmt.Errorf("invalid %s_before: %w", name, err)
	}

	if start != nil && end != nil && !start.before(*end) {
		return dailyWindow{}, fmt.Errorf("%s_after must be before %s_before", name, name)
	}

	return dailyWindow{
		start: start,
		end:   end,
	}, nil
}

func parseClockTime(value string) (*clockTime, error) {
	if value == "" {
		return nil, nil
	}

	parsed, err := time.Parse("15:04", value)
	if err != nil {
		return nil, fmt.Errorf("expected HH:MM, got %q", value)
	}

	return &clockTime{hour: parsed.Hour(), minute: parsed.Minute()}, nil
}

func clampToDailyWindow(runAt, windowDate time.Time, window dailyWindow) time.Time {
	if window.start != nil {
		start := window.start.onDate(windowDate)
		if runAt.Before(start) {
			return start
		}
	}

	if window.end != nil {
		end := window.end.onDate(windowDate)
		if runAt.After(end) {
			return end
		}
	}

	return runAt
}
