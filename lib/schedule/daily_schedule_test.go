package schedule

import (
	"reflect"
	"testing"
)

func TestDailyScheduleAt(t *testing.T) {
	ds := NewDailySchedule()
	slot1 := &TimeSlot{Start: NewTime(8, 0, 0), Stop: NewTimePointer(10, 0, 0)}
	slot2 := &TimeSlot{Start: NewTime(10, 30, 0)}
	slot3 := &TimeSlot{Start: NewTime(11, 0, 0)}
	slot4 := &TimeSlot{Start: NewTime(14, 0, 0), Stop: NewTimePointer(16, 0, 0)}
	slot5 := &TimeSlot{Start: NewTime(20, 0, 0)}

	for _, slot := range []*TimeSlot{slot1, slot2, slot3, slot4, slot5} {
		err := ds.Add(slot)
		if err != nil {
			t.Fatal(err)
		}
	}

	tt := []struct {
		name     string
		at       Time
		expected *TimeSlot
	}{
		{name: "expect no slot found", at: NewTime(1, 0, 0), expected: nil},
		{name: "expect slot1", at: NewTime(9, 0, 0), expected: slot1},
		{name: "expect slot 2", at: NewTime(10, 35, 0), expected: slot2},
		{name: "expect not slot found again", at: NewTime(10, 25, 0), expected: nil},
		{name: "expect slot 3", at: NewTime(13, 0, 0), expected: slot3},
		{name: "expect slot 5", at: NewTime(22, 0, 0), expected: slot5},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			slot := ds.At(tc.at)
			if !reflect.DeepEqual(slot, tc.expected) {
				t.Errorf("expected %+v, got %+v", tc.expected, slot)
			}
		})
	}
}

func TestDailyScheduleAdd(t *testing.T) {
	ds := NewDailySchedule()

	tt := []struct {
		name        string
		newTimeSlot *TimeSlot
		expected    error
	}{
		{name: "nil timeslot"},
		{
			name:        "first timeslot",
			newTimeSlot: &TimeSlot{Start: NewTime(9, 0, 0)},
		},
		{
			name:        "exact same timeslot",
			newTimeSlot: &TimeSlot{Start: NewTime(9, 0, 0)},
			expected:    ErrOverlappingTimeslots,
		},
		{
			name: "second timeslot",
			newTimeSlot: &TimeSlot{
				Start: NewTime(10, 0, 0),
				Stop:  NewTimePointer(11, 0, 0),
			},
		},
		{
			name: "in previous timeslot",
			newTimeSlot: &TimeSlot{
				Start: NewTime(10, 10, 0),
				Stop:  NewTimePointer(10, 20, 0),
			},
			expected: ErrOverlappingTimeslots,
		},
		{
			name: "new timeslot starting at the end of the previous one",
			newTimeSlot: &TimeSlot{
				Start: NewTime(11, 00, 0),
			},
		},
		{
			name: "overlapping with another timeslot",
			newTimeSlot: &TimeSlot{
				Start: NewTime(9, 30, 0),
				Stop:  NewTimePointer(10, 30, 0),
			},
			expected: ErrOverlappingTimeslots,
		},
		{
			name:        "before everything else",
			newTimeSlot: &TimeSlot{Start: NewTime(0, 0, 0)},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got := ds.Add(tc.newTimeSlot)
			if got != tc.expected {
				t.Errorf("expected %+v, got %+v", tc.expected, got)
			}
		})
	}
}

func TestDailyScheduleDelete(t *testing.T) {
	ds := NewDailySchedule()

	fakeUUID := "254f1bff-a6a8-422f-889f-99195664ca77"
	slot1 := &TimeSlot{Start: NewTime(8, 0, 0)}
	slot2 := &TimeSlot{Start: NewTime(9, 0, 0), UUID: fakeUUID}
	slot3 := &TimeSlot{Start: NewTime(10, 0, 0)}

	for _, ts := range []*TimeSlot{slot1, slot2, slot3} {
		err := ds.Add(ts)
		if err != nil {
			t.Fatal(err)
		}
	}

	if slot3.UUID == "" {
		t.Errorf("slot3 should have a UUID by now")
	}

	if err := ds.Delete(fakeUUID); err != nil {
		t.Fatal(err)
	}

	if ds.Len() != 2 {
		t.Fatalf("invalid schedule size")
	}

	if ds[0] != slot1 {
		t.Fatalf("expected slot 1")
	}
	if ds[1] != slot3 {
		t.Fatalf("expected slot 3")
	}

	if err := ds.Delete(slot3.UUID); err != nil {
		t.Fatal(err)
	}

	if err := ds.Delete(slot1.UUID); err != nil {
		t.Fatal(err)
	}

	if err := ds.Delete("invalid id"); err != ErrTimeSlotNotFound {
		t.Error("this invalid id should not work")
	}
}
