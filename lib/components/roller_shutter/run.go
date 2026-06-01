package rollershutter

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/gregdel/homed/lib/components"
)

const failedActionRetryDelay = 30 * time.Second

// Run implements the component interface
func (rs *RollerShutter) Run(ctx context.Context, logger *slog.Logger, _ *components.Components) error {
	if err := rs.YAMLParams.Decode(&rs.Params); err != nil {
		return err
	}
	if err := rs.parseDailyWindows(); err != nil {
		return err
	}

	log := logger.With(slog.String("component_id", rs.ID()))
	rs.logger = log

	if !rs.Params.Enabled {
		log.Info("automatic control disabled")
		return nil
	}

	nextAction := rs.nextAction()
	var retryAt time.Time

	for {
		runAt := rs.nextRunTime(nextAction)
		if !retryAt.IsZero() && retryAt.After(runAt) {
			runAt = retryAt
		}
		sleepDuration := time.Until(runAt)

		log.Info("setting next roller shutter action",
			slog.Time("run_at", runAt),
			slog.String("action", nextAction.String()),
			slog.Duration("sleep_duration", sleepDuration),
		)

		select {
		case <-ctx.Done():
			log.Info("stopping")
			return nil
		case <-time.After(sleepDuration):
			completedAction := nextAction
			err := rs.performAction(completedAction)

			if err != nil {
				log.Info("failed to run roller shutter action",
					slog.String("action", completedAction.String()),
					slog.Any("error", err),
				)
			}

			nextAction, retryAt = rs.nextActionAfterAttemptAt(time.Now(), completedAction, err)
			if err != nil && !retryAt.IsZero() {
				log.Info("retrying roller shutter action",
					slog.String("action", nextAction.String()),
					slog.Time("retry_at", retryAt),
				)
			}
		}
	}
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

func (action shutterAction) followingAction() shutterAction {
	switch action {
	case openShutters:
		return closeShutters
	case closeShutters:
		return openShutters
	default:
		return openShutters
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

func (rs *RollerShutter) nextActionAfterAttemptAt(now time.Time, action shutterAction, err error) (shutterAction, time.Time) {
	if err == nil {
		return action.followingAction(), time.Time{}
	}

	retryAt, ok := rs.retryRunTimeAt(now, action)
	if ok {
		return action, retryAt
	}

	return action.followingAction(), time.Time{}
}

func (rs *RollerShutter) retryRunTimeAt(now time.Time, action shutterAction) (time.Time, bool) {
	runAt, deadline := rs.runWindowForActionAt(now, action)
	if !isWithinCatchUpWindow(now, runAt, deadline) {
		return time.Time{}, false
	}

	retryAt := now.Add(failedActionRetryDelay)
	if retryAt.After(deadline) {
		retryAt = deadline
	}
	if !retryAt.After(now) {
		return time.Time{}, false
	}

	return retryAt, true
}

func (rs *RollerShutter) nextAction() shutterAction {
	now := time.Now()
	return rs.nextActionAt(now)
}

func (rs *RollerShutter) nextActionAt(now time.Time) shutterAction {
	openAt, openDeadline := rs.runWindowForActionAt(now, openShutters)
	if now.Before(openAt) || isWithinCatchUpWindow(now, openAt, openDeadline) {
		return openShutters
	}

	closeAt, closeDeadline := rs.runWindowForActionAt(now, closeShutters)
	if now.Before(closeAt) || isWithinCatchUpWindow(now, closeAt, closeDeadline) {
		return closeShutters
	}

	return openShutters
}

func (rs *RollerShutter) nextRunTime(action shutterAction) time.Time {
	now := time.Now()
	return rs.nextRunTimeAt(now, action)
}

func (rs *RollerShutter) nextRunTimeAt(now time.Time, action shutterAction) time.Time {
	runAt, deadline := rs.runWindowForActionAt(now, action)
	if now.Before(runAt) {
		return runAt
	}

	if isWithinCatchUpWindow(now, runAt, deadline) {
		return now
	}

	nextRunAt, _ := rs.runWindowForActionAt(now.AddDate(0, 0, 1), action)
	return nextRunAt
}

func (rs *RollerShutter) runWindowForActionAt(now time.Time, action shutterAction) (runAt, deadline time.Time) {
	switch action {
	case openShutters:
		runAt := rs.scheduledOpenTimeAt(now)
		return runAt, latestRunTime(now, runAt, rs.openWindow)
	case closeShutters:
		runAt := rs.scheduledCloseTimeAt(now)
		return runAt, latestRunTime(now, runAt, rs.closeWindow)
	default:
		return time.Time{}, time.Time{}
	}
}

func (rs *RollerShutter) scheduledOpenTimeAt(now time.Time) time.Time {
	sunrise, _ := rs.getSunriseSunset(now)

	return moveInsideDailyWindow(sunrise.In(now.Location()), now, rs.openWindow)
}

func (rs *RollerShutter) scheduledCloseTimeAt(now time.Time) time.Time {
	_, sunset := rs.getSunriseSunset(now)

	return moveInsideDailyWindow(sunset.In(now.Location()), now, rs.closeWindow)
}

func isWithinCatchUpWindow(now, runAt, deadline time.Time) bool {
	return !now.Before(runAt) && !now.After(deadline)
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
	earliest    time.Duration
	hasEarliest bool
	latest      time.Duration
	hasLatest   bool
}

func parseDailyWindow(name, afterValue, beforeValue string) (dailyWindow, error) {
	earliest, hasEarliest, err := parseClockTime(afterValue)
	if err != nil {
		return dailyWindow{}, fmt.Errorf("invalid %s_after: %w", name, err)
	}

	latest, hasLatest, err := parseClockTime(beforeValue)
	if err != nil {
		return dailyWindow{}, fmt.Errorf("invalid %s_before: %w", name, err)
	}

	if hasEarliest && hasLatest && earliest >= latest {
		return dailyWindow{}, fmt.Errorf("%s_after must be before %s_before", name, name)
	}

	return dailyWindow{
		earliest:    earliest,
		hasEarliest: hasEarliest,
		latest:      latest,
		hasLatest:   hasLatest,
	}, nil
}

func parseClockTime(value string) (time.Duration, bool, error) {
	if value == "" {
		return 0, false, nil
	}

	parsed, err := time.Parse("15:04", value)
	if err != nil {
		return 0, false, fmt.Errorf("expected HH:MM, got %q", value)
	}

	return time.Duration(parsed.Hour())*time.Hour + time.Duration(parsed.Minute())*time.Minute, true, nil
}

func moveInsideDailyWindow(runAt, windowDate time.Time, window dailyWindow) time.Time {
	if window.hasEarliest {
		earliest := clockTimeOnDate(windowDate, window.earliest)
		if runAt.Before(earliest) {
			return earliest
		}
	}

	if window.hasLatest {
		latest := clockTimeOnDate(windowDate, window.latest)
		if runAt.After(latest) {
			return latest
		}
	}

	return runAt
}

func latestRunTime(windowDate, runAt time.Time, window dailyWindow) time.Time {
	if window.hasLatest {
		return clockTimeOnDate(windowDate, window.latest)
	}

	return runAt
}

func clockTimeOnDate(date time.Time, clock time.Duration) time.Time {
	midnight := time.Date(
		date.Year(),
		date.Month(),
		date.Day(),
		0,
		0,
		0,
		0,
		date.Location(),
	)

	return midnight.Add(clock)
}
