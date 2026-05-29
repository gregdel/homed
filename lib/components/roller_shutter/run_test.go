package rollershutter

import (
	"errors"
	"testing"
	"time"

	"github.com/gregdel/homed/lib/config"
)

func TestScheduledOpenTimeAt(t *testing.T) {
	tt := []struct {
		name   string
		now    time.Time
		params Params
		want   func(rs *RollerShutter, now time.Time) time.Time
	}{
		{
			name: "sunrise before open_after clamps to open_after",
			now:  time.Date(2026, 5, 29, 12, 0, 0, 0, time.FixedZone("test", 2*60*60)),
			params: Params{
				Location:   testLocation(),
				OpenAfter:  "07:00",
				OpenBefore: "09:00",
			},
			want: func(_ *RollerShutter, now time.Time) time.Time {
				return time.Date(2026, 5, 29, 7, 0, 0, 0, now.Location())
			},
		},
		{
			name: "sunrise inside window uses sunrise",
			now:  time.Date(2026, 12, 29, 12, 0, 0, 0, time.FixedZone("test", 1*60*60)),
			params: Params{
				Location:   testLocation(),
				OpenAfter:  "07:00",
				OpenBefore: "09:00",
			},
			want: func(rs *RollerShutter, now time.Time) time.Time {
				sunrise, _ := rs.getSunriseSunset(now)
				return sunrise.In(now.Location())
			},
		},
		{
			name: "sunrise after open_before clamps to open_before",
			now:  time.Date(2026, 12, 29, 12, 0, 0, 0, time.FixedZone("test", 1*60*60)),
			params: Params{
				Location:   testLocation(),
				OpenAfter:  "06:00",
				OpenBefore: "07:00",
			},
			want: func(_ *RollerShutter, now time.Time) time.Time {
				return time.Date(2026, 12, 29, 7, 0, 0, 0, now.Location())
			},
		},
		{
			name: "open_after without open_before clamps early sunrise",
			now:  time.Date(2026, 5, 29, 1, 0, 0, 0, time.FixedZone("CEST", 2*60*60)),
			params: Params{
				Location:  testLocation(),
				OpenAfter: "07:00",
			},
			want: func(_ *RollerShutter, now time.Time) time.Time {
				return time.Date(2026, 5, 29, 7, 0, 0, 0, now.Location())
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			rs := testRollerShutter(t, tc.params)
			assertTimeEqual(t, rs.scheduledOpenTimeAt(tc.now), tc.want(rs, tc.now))
		})
	}
}

func TestScheduledCloseTimeAt(t *testing.T) {
	tt := []struct {
		name   string
		now    time.Time
		params Params
		want   func(rs *RollerShutter, now time.Time) time.Time
	}{
		{
			name: "sunset after close_before clamps to close_before",
			now:  time.Date(2026, 5, 29, 12, 0, 0, 0, time.FixedZone("test", 2*60*60)),
			params: Params{
				Location:    testLocation(),
				CloseAfter:  "18:00",
				CloseBefore: "21:30",
			},
			want: func(_ *RollerShutter, now time.Time) time.Time {
				return time.Date(2026, 5, 29, 21, 30, 0, 0, now.Location())
			},
		},
		{
			name: "sunset inside window uses sunset",
			now:  time.Date(2026, 5, 29, 12, 0, 0, 0, time.FixedZone("test", 2*60*60)),
			params: Params{
				Location:    testLocation(),
				CloseAfter:  "18:00",
				CloseBefore: "23:00",
			},
			want: func(rs *RollerShutter, now time.Time) time.Time {
				_, sunset := rs.getSunriseSunset(now)
				return sunset.In(now.Location())
			},
		},
		{
			name: "sunset before close_after clamps to close_after",
			now:  time.Date(2026, 12, 29, 12, 0, 0, 0, time.FixedZone("test", 1*60*60)),
			params: Params{
				Location:    testLocation(),
				CloseAfter:  "18:00",
				CloseBefore: "22:30",
			},
			want: func(_ *RollerShutter, now time.Time) time.Time {
				return time.Date(2026, 12, 29, 18, 0, 0, 0, now.Location())
			},
		},
		{
			name: "close_after without close_before clamps early sunset",
			now:  time.Date(2026, 12, 29, 12, 0, 0, 0, time.FixedZone("CET", 1*60*60)),
			params: Params{
				Location:   testLocation(),
				CloseAfter: "18:00",
			},
			want: func(_ *RollerShutter, now time.Time) time.Time {
				return time.Date(2026, 12, 29, 18, 0, 0, 0, now.Location())
			},
		},
		{
			name: "sunset after local midnight clamps to schedule date",
			now:  time.Date(2026, 6, 21, 12, 0, 0, 0, time.FixedZone("UTC+14", 14*60*60)),
			params: Params{
				Location:    testLocation(),
				CloseAfter:  "18:00",
				CloseBefore: "22:30",
			},
			want: func(_ *RollerShutter, now time.Time) time.Time {
				return time.Date(2026, 6, 21, 22, 30, 0, 0, now.Location())
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			rs := testRollerShutter(t, tc.params)
			assertTimeEqual(t, rs.scheduledCloseTimeAt(tc.now), tc.want(rs, tc.now))
		})
	}
}

func TestParseDailyWindows(t *testing.T) {
	tt := []struct {
		name    string
		params  Params
		wantErr bool
	}{
		{
			name: "valid windows",
			params: Params{
				OpenAfter:   "07:00",
				OpenBefore:  "09:00",
				CloseAfter:  "18:00",
				CloseBefore: "22:30",
			},
		},
		{
			name: "invalid format",
			params: Params{
				OpenAfter: "7",
			},
			wantErr: true,
		},
		{
			name: "invalid ordering",
			params: Params{
				OpenAfter:  "09:00",
				OpenBefore: "07:00",
			},
			wantErr: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			rs := &RollerShutter{Params: tc.params}
			err := rs.parseDailyWindows()
			if tc.wantErr && err == nil {
				t.Fatal("expected error")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
		})
	}
}

func TestNextActionAfterAttemptAt(t *testing.T) {
	errFailed := errors.New("failed")
	standardParams := Params{
		Location:    testLocation(),
		OpenAfter:   "07:00",
		OpenBefore:  "09:00",
		CloseAfter:  "18:00",
		CloseBefore: "22:30",
	}

	tt := []struct {
		name            string
		params          Params
		now             time.Time
		action          shutterAction
		err             error
		wantAction      shutterAction
		wantRetryAt     time.Time
		wantNoRetryTime bool
	}{
		{
			name:       "advances after successful action",
			params:     standardParams,
			now:        time.Date(2026, 5, 29, 7, 30, 0, 0, time.FixedZone("CEST", 2*60*60)),
			action:     openShutters,
			wantAction: closeShutters,
		},
		{
			name:        "retries failed open before deadline",
			params:      standardParams,
			now:         time.Date(2026, 5, 29, 7, 30, 0, 0, time.FixedZone("CEST", 2*60*60)),
			action:      openShutters,
			err:         errFailed,
			wantAction:  openShutters,
			wantRetryAt: time.Date(2026, 5, 29, 7, 30, 30, 0, time.FixedZone("CEST", 2*60*60)),
		},
		{
			name:        "caps retry at deadline",
			params:      standardParams,
			now:         time.Date(2026, 5, 29, 8, 59, 45, 0, time.FixedZone("CEST", 2*60*60)),
			action:      openShutters,
			err:         errFailed,
			wantAction:  openShutters,
			wantRetryAt: time.Date(2026, 5, 29, 9, 0, 0, 0, time.FixedZone("CEST", 2*60*60)),
		},
		{
			name:            "advances failed open after deadline",
			params:          standardParams,
			now:             time.Date(2026, 5, 29, 9, 0, 1, 0, time.FixedZone("CEST", 2*60*60)),
			action:          openShutters,
			err:             errFailed,
			wantAction:      closeShutters,
			wantNoRetryTime: true,
		},
		{
			name: "retries failed close before deadline",
			params: Params{
				Location:    testLocation(),
				CloseAfter:  "18:00",
				CloseBefore: "23:00",
			},
			now:         time.Date(2026, 5, 29, 22, 45, 0, 0, time.FixedZone("CEST", 2*60*60)),
			action:      closeShutters,
			err:         errFailed,
			wantAction:  closeShutters,
			wantRetryAt: time.Date(2026, 5, 29, 22, 45, 30, 0, time.FixedZone("CEST", 2*60*60)),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			rs := testRollerShutter(t, tc.params)
			gotAction, gotRetryAt := rs.nextActionAfterAttemptAt(tc.now, tc.action, tc.err)
			if gotAction != tc.wantAction {
				t.Fatalf("expected action %s, got %s", tc.wantAction, gotAction)
			}
			if tc.wantNoRetryTime || tc.wantRetryAt.IsZero() {
				if !gotRetryAt.IsZero() {
					t.Fatalf("expected no retry time, got %s", gotRetryAt)
				}
				return
			}

			assertTimeEqual(t, gotRetryAt, tc.wantRetryAt)
		})
	}
}

func TestNextActionAndRunTimeAt(t *testing.T) {
	standardParams := Params{
		Location:    testLocation(),
		OpenAfter:   "07:00",
		OpenBefore:  "09:00",
		CloseAfter:  "18:00",
		CloseBefore: "22:30",
	}

	tt := []struct {
		name       string
		params     Params
		now        time.Time
		action     shutterAction
		wantAction shutterAction
		wantRunAt  func(rs *RollerShutter, now time.Time) time.Time
	}{
		{
			name:       "schedules normal morning open",
			params:     standardParams,
			now:        time.Date(2026, 5, 29, 6, 30, 0, 0, time.FixedZone("CEST", 2*60*60)),
			action:     openShutters,
			wantAction: openShutters,
			wantRunAt: func(_ *RollerShutter, now time.Time) time.Time {
				return time.Date(2026, 5, 29, 7, 0, 0, 0, now.Location())
			},
		},
		{
			name:       "catches up morning open before deadline",
			params:     standardParams,
			now:        time.Date(2026, 5, 29, 7, 30, 0, 0, time.FixedZone("CEST", 2*60*60)),
			action:     openShutters,
			wantAction: openShutters,
			wantRunAt: func(_ *RollerShutter, now time.Time) time.Time {
				return now
			},
		},
		{
			name:       "skips late morning open after deadline",
			params:     standardParams,
			now:        time.Date(2026, 5, 29, 9, 30, 0, 0, time.FixedZone("CEST", 2*60*60)),
			action:     closeShutters,
			wantAction: closeShutters,
			wantRunAt: func(rs *RollerShutter, now time.Time) time.Time {
				return rs.scheduledCloseTimeAt(now)
			},
		},
		{
			name: "does not catch up morning open without deadline",
			params: Params{
				Location:    testLocation(),
				OpenAfter:   "07:00",
				CloseAfter:  "18:00",
				CloseBefore: "22:30",
			},
			now:        time.Date(2026, 5, 29, 7, 30, 0, 0, time.FixedZone("CEST", 2*60*60)),
			action:     openShutters,
			wantAction: closeShutters,
			wantRunAt: func(_ *RollerShutter, now time.Time) time.Time {
				nextDay := now.AddDate(0, 0, 1)
				return time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(), 7, 0, 0, 0, now.Location())
			},
		},
		{
			name: "catches up evening close within deadline",
			params: Params{
				Location:    testLocation(),
				CloseAfter:  "18:00",
				CloseBefore: "23:00",
			},
			now:        time.Date(2026, 5, 29, 22, 45, 0, 0, time.FixedZone("CEST", 2*60*60)),
			action:     closeShutters,
			wantAction: closeShutters,
			wantRunAt: func(_ *RollerShutter, now time.Time) time.Time {
				return now
			},
		},
		{
			name: "does not catch up after close_before",
			params: Params{
				Location:    testLocation(),
				CloseAfter:  "18:00",
				CloseBefore: "21:30",
			},
			now:        time.Date(2026, 5, 29, 21, 31, 0, 0, time.FixedZone("CEST", 2*60*60)),
			action:     openShutters,
			wantAction: openShutters,
			wantRunAt: func(rs *RollerShutter, now time.Time) time.Time {
				return rs.scheduledOpenTimeAt(now.AddDate(0, 0, 1))
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			rs := testRollerShutter(t, tc.params)
			if got := rs.nextActionAt(tc.now); got != tc.wantAction {
				t.Fatalf("expected nextActionAt to be %s, got %s", tc.wantAction, got)
			}

			got := rs.nextRunTimeAt(tc.now, tc.action)
			assertTimeEqual(t, got, tc.wantRunAt(rs, tc.now))
		})
	}
}

func testRollerShutter(t *testing.T, params Params) *RollerShutter {
	t.Helper()

	rs := &RollerShutter{Params: params}
	if err := rs.parseDailyWindows(); err != nil {
		t.Fatal(err)
	}

	return rs
}

func testLocation() config.Location {
	return config.Location{
		Latitude:  50.62955091282183,
		Longitude: 3.055823770700181,
	}
}

func assertTimeEqual(t *testing.T, got, expected time.Time) {
	t.Helper()

	if !got.Equal(expected) {
		t.Fatalf("expected %s, got %s", expected, got)
	}
}
