package schedule

import (
	"testing"
	"time"
)

func TestSchedule(t *testing.T) {
	// Fake the now function
	setNow(NewTimePointer(11, 00, 00))
	defer setNow(nil)

	schedule := New(10)

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

	if err := schedule.Delete(day, slot1.ID); err != nil {
		t.Fatalf("expected to be able to delete a slot, got %s", err.Error())
	}
}

func TestScheduleValue(t *testing.T) {
	// Fake the now function
	setNow(NewTimePointer(11, 00, 00))
	defer setNow(nil)

	var expected float64 = 10
	schedule := New(expected)

	got := schedule.Value()
	if got != expected {
		t.Errorf("expected %f, got %f", expected, got)
	}

	expected = 20
	day := now().Weekday()
	ts := &TimeSlot{
		Start: NewTime(0, 0, 0),
		Stop:  NewTimePointer(12, 0, 0),
		Value: expected,
	}

	if err := schedule.Add(day, ts); err != nil {
		t.Fatal(err)
	}

	got = schedule.Value()
	if got != expected {
		t.Errorf("expected %f, got %f", expected, got)
	}

	expected = 18
	o := NewOverride(
		now().Add(-1*time.Hour),
		now().Add(1*time.Hour),
		expected,
	)
	if err := schedule.AddOverride(o); err != nil {
		t.Fatal(err)
	}

	got = schedule.Value()
	if got != expected {
		t.Errorf("expected %f, got %f", expected, got)
	}

	if err := schedule.DeleteOverride(o.ID); err != nil {
		t.Fatal(err)
	}

	expected = 20
	got = schedule.Value()
	if got != expected {
		t.Errorf("expected %f, got %f", expected, got)
	}
}

func TestScheduleNextChange(t *testing.T) {
	// Fake the now function
	setNow(NewTimePointer(11, 00, 00))
	defer setNow(nil)

	today := now().Weekday()
	tomorrow := now().Add(1 * 24 * time.Hour).Weekday()

	todaySlot1 := &TimeSlot{Start: NewTime(16, 0, 0), Stop: NewTimePointer(17, 0, 0)}
	todaySlot2 := &TimeSlot{Start: NewTime(18, 0, 0)}
	tomorrowSlot1 := &TimeSlot{Start: NewTime(9, 0, 0), Stop: NewTimePointer(12, 0, 0)}

	t1 := now().Add(5 * time.Hour)
	t2 := now().Add(6 * time.Hour)
	t3 := now().Add(7 * time.Hour)
	t4 := now().Add(22 * time.Hour)

	schedule := New(10)
	emptyChange := schedule.NextChange()
	if emptyChange != nil {
		t.Errorf("next change should be nil")
	}

	for _, todo := range []struct {
		wd time.Weekday
		ts *TimeSlot
	}{
		{wd: today, ts: todaySlot1},
		{wd: today, ts: todaySlot2},
		{wd: tomorrow, ts: tomorrowSlot1},
	} {
		if err := schedule.Add(todo.wd, todo.ts); err != nil {
			t.Fatal(err)
		}
	}

	data := []struct {
		name     string
		wd       time.Weekday
		when     *Time
		expected *time.Time
	}{
		{name: "before", wd: today, when: NewTimePointer(11, 30, 0), expected: &t1},
		{name: "in timeslot", wd: today, when: NewTimePointer(16, 30, 0), expected: &t2},
		{name: "between timeslots", wd: today, when: NewTimePointer(17, 30, 0), expected: &t3},
		{name: "without stop", wd: today, when: NewTimePointer(18, 30, 0), expected: &t4},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			setNow(d.when)

			nextChange := schedule.NextChange()
			if nextChange == nil {
				t.Fatalf("expected next change, got nothing")
			}

			if *nextChange != *d.expected {
				t.Fatalf("expected %s, got %s", d.expected, nextChange)
			}
		})
	}
}
