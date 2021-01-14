package schedule

import (
	"testing"
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
