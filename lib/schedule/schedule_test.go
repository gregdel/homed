package schedule

import (
	"testing"
	"time"
)

func TestSchedule(t *testing.T) {
	// Fake the now function
	setNow(NewTimePointer(11, 00, 00))
	defer setNow(nil)

	schedule := New()

	day := now().Weekday()
	slot1 := &TimeSlot{Start: NewTime(0, 0, 0)}
	slot2 := &TimeSlot{Start: NewTime(10, 0, 0)}
	slot3 := &TimeSlot{Start: NewTime(13, 0, 0)}

	for _, slot := range []*TimeSlot{slot1, slot2, slot3} {
		if err := schedule.Add(day, slot); err != nil {
			t.Fatal(err)
		}
	}

	got := schedule.Now()
	if got == nil || *got != *slot2 {
		t.Errorf("expected to get slot2, got %+v", got)
	}

	if err := schedule.Delete(day, slot1.UUID); err != nil {
		t.Fatalf("expected to be able to delete a slot, got %s", err.Error())
	}
}

func TestScheduleNext(t *testing.T) {
	// Fake the now function
	setNow(NewTimePointer(11, 00, 00))
	defer setNow(nil)

	today := now().Weekday()
	tomorrow := now().Add(1 * 24 * time.Hour).Weekday()
	inTwoDays := now().Add(2 * 24 * time.Hour).Weekday()
	inSixDays := now().Add(6 * 24 * time.Hour).Weekday()

	slot1 := &TimeSlot{Start: NewTime(16, 0, 0)}
	slot2 := &TimeSlot{Start: NewTime(10, 30, 40)}
	slot3 := &TimeSlot{Start: NewTime(13, 0, 0)}

	t1 := now().Add(6 * 24 * time.Hour).Add(2 * time.Hour)
	t2 := now().Add(2 * 24 * time.Hour).Add(2 * time.Hour)
	t3 := now().Add(24 * time.Hour).Add(-29 * time.Minute).Add(-20 * time.Second)
	t4 := now().Add(5 * time.Hour)

	data := []struct {
		name string
		wd   time.Weekday
		ts   *TimeSlot
		nt   *time.Time
	}{
		{name: "nothing", wd: today, ts: nil},
		{name: "in six days", wd: inSixDays, ts: slot3, nt: &t1},
		{name: "in two days", wd: inTwoDays, ts: slot3, nt: &t2},
		{name: "tomorrow", wd: tomorrow, ts: slot2, nt: &t3},
		{name: "today", wd: today, ts: slot1, nt: &t4},
	}

	schedule := New()
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			if err := schedule.Add(d.wd, d.ts); err != nil {
				t.Fatal(err)
			}

			got, gotWeekday := schedule.Next()
			if d.ts == nil {
				if got != nil {
					t.Errorf("expected nothing got %+v", got)
				}

				_, gotTime := schedule.NextTime()
				if gotTime != nil {
					t.Errorf("expected not time got %+v", gotTime)
				}

				return
			}

			if got == nil {
				t.Fatalf("expected something, got nothing")
			}

			if *got != *d.ts {
				t.Errorf("expected to get %+v, got %+v", d.ts, got)
			}

			if gotWeekday != d.wd {
				t.Errorf("expected %s, got %s", d.wd.String(), gotWeekday.String())
			}

			if d.nt != nil {
				_, gotTime := schedule.NextTime()
				if gotTime == nil {
					t.Errorf("expected a time got nil")
					return
				}

				if *gotTime != *d.nt {
					t.Errorf("expected time %s, got %s", d.nt, gotTime)
				}
			}
		})
	}
}
