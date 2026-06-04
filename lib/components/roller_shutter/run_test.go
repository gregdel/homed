package rollershutter

import (
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
			name: "midnight is a valid configured window bound",
			params: Params{
				OpenAfter:  "00:00",
				OpenBefore: "07:00",
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

func TestNextScheduledEventAt(t *testing.T) {
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
		wantAction shutterAction
		wantRunAt  func(rs *RollerShutter, now time.Time) time.Time
	}{
		{
			name:       "schedules normal morning open",
			params:     standardParams,
			now:        time.Date(2026, 5, 29, 6, 30, 0, 0, time.FixedZone("CEST", 2*60*60)),
			wantAction: openShutters,
			wantRunAt: func(_ *RollerShutter, now time.Time) time.Time {
				return time.Date(2026, 5, 29, 7, 0, 0, 0, now.Location())
			},
		},
		{
			name:       "runs open at exact open time",
			params:     standardParams,
			now:        time.Date(2026, 5, 29, 7, 0, 0, 0, time.FixedZone("CEST", 2*60*60)),
			wantAction: openShutters,
			wantRunAt: func(_ *RollerShutter, now time.Time) time.Time {
				return now
			},
		},
		{
			name:       "skips missed open before open_before",
			params:     standardParams,
			now:        time.Date(2026, 5, 29, 7, 30, 0, 0, time.FixedZone("CEST", 2*60*60)),
			wantAction: closeShutters,
			wantRunAt: func(rs *RollerShutter, now time.Time) time.Time {
				return rs.scheduledCloseTimeAt(now)
			},
		},
		{
			name:       "after open before close schedules close",
			params:     standardParams,
			now:        time.Date(2026, 5, 29, 9, 30, 0, 0, time.FixedZone("CEST", 2*60*60)),
			wantAction: closeShutters,
			wantRunAt: func(rs *RollerShutter, now time.Time) time.Time {
				return rs.scheduledCloseTimeAt(now)
			},
		},
		{
			name: "does not catch up open without open_before",
			params: Params{
				Location:    testLocation(),
				OpenAfter:   "07:00",
				CloseAfter:  "18:00",
				CloseBefore: "22:30",
			},
			now:        time.Date(2026, 5, 29, 7, 30, 0, 0, time.FixedZone("CEST", 2*60*60)),
			wantAction: closeShutters,
			wantRunAt: func(rs *RollerShutter, now time.Time) time.Time {
				return rs.scheduledCloseTimeAt(now)
			},
		},
		{
			name: "runs close at exact close time",
			params: Params{
				Location:    testLocation(),
				CloseAfter:  "18:00",
				CloseBefore: "23:00",
			},
			now:        time.Date(2026, 12, 29, 18, 0, 0, 0, time.FixedZone("CET", 1*60*60)),
			wantAction: closeShutters,
			wantRunAt: func(rs *RollerShutter, now time.Time) time.Time {
				return rs.scheduledCloseTimeAt(now)
			},
		},
		{
			name: "skips missed close before close_before",
			params: Params{
				Location:    testLocation(),
				CloseAfter:  "18:00",
				CloseBefore: "23:00",
			},
			now:        time.Date(2026, 5, 29, 22, 45, 0, 0, time.FixedZone("CEST", 2*60*60)),
			wantAction: openShutters,
			wantRunAt: func(rs *RollerShutter, now time.Time) time.Time {
				return rs.scheduledOpenTimeAt(now.AddDate(0, 0, 1))
			},
		},
		{
			name: "does not catch up close without close_before",
			params: Params{
				Location:   testLocation(),
				OpenAfter:  "07:00",
				CloseAfter: "18:00",
			},
			now:        time.Date(2026, 12, 29, 19, 0, 0, 0, time.FixedZone("CET", 1*60*60)),
			wantAction: openShutters,
			wantRunAt: func(rs *RollerShutter, now time.Time) time.Time {
				return rs.scheduledOpenTimeAt(now.AddDate(0, 0, 1))
			},
		},
		{
			name: "after close schedules tomorrow open",
			params: Params{
				Location:    testLocation(),
				CloseAfter:  "18:00",
				CloseBefore: "21:30",
			},
			now:        time.Date(2026, 5, 29, 21, 31, 0, 0, time.FixedZone("CEST", 2*60*60)),
			wantAction: openShutters,
			wantRunAt: func(rs *RollerShutter, now time.Time) time.Time {
				return rs.scheduledOpenTimeAt(now.AddDate(0, 0, 1))
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			rs := testRollerShutter(t, tc.params)

			gotEvent := rs.nextEvent(tc.now)
			if gotEvent.action != tc.wantAction {
				t.Fatalf("expected action %s, got %s", tc.wantAction, gotEvent.action)
			}

			assertTimeEqual(t, gotEvent.scheduledAt, tc.wantRunAt(rs, tc.now))
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
