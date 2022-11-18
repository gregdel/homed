package schedule

import (
	"errors"
	"sort"
)

var (
	// ErrOverlappingTimeslots is returns two timeslots overlaps
	ErrOverlappingTimeslots = errors.New("schedule: overlapping timeslots")
	// ErrStopBeforeStart is returns if the stop is before the start
	ErrStopBeforeStart = errors.New("schedule: stop before start")
	// ErrTimeSlotNotFound is returned if the given timeslot is not found
	ErrTimeSlotNotFound = errors.New("schedule: timeslot not found")
)

// DailySchedule holds a sorted slice of time slots
type DailySchedule []*TimeSlot

// NewDailySchedule returns a new DailySchedule
func NewDailySchedule() DailySchedule {
	return []*TimeSlot{}
}

// Len implements the sort interface
func (ds DailySchedule) Len() int {
	return len(ds)
}

// Less implements the sort interface
func (ds DailySchedule) Less(i, j int) bool {
	return ds[i].Start.Before(ds[j].Start)
}

// Swap implements the sort interface
func (ds DailySchedule) Swap(i, j int) {
	ds[i], ds[j] = ds[j], ds[i]
}

// Add adds a time slot to the daily schedule
func (ds *DailySchedule) Add(n *TimeSlot) error {
	if n == nil {
		return nil
	}

	if err := n.validate(); err != nil {
		return err
	}

	// Check if there is a timeslot at the new time slot start
	s := ds.At(n.Start)
	if s != nil && ((s.Start == n.Start) || s.Stop != nil) {
		// We're in the middle of the timeslot
		return ErrOverlappingTimeslots
	}

	// Get the fist timeslot after the start of the new timeslot
	s = ds.FirstAfter(n.Start)
	if n.Stop != nil && s != nil && n.Stop.After(s.Start) {
		// The end of the new time slot is after the beginning the previous one
		return ErrOverlappingTimeslots
	}

	n.generateID()
	*ds = append(*ds, n)

	// Sort the schedule
	sort.Sort(ds)
	return nil
}

// Delete deletes a timeslot
func (ds *DailySchedule) Delete(id string) error {
	var found *int
	for i, ts := range *ds {
		if ts.ID == id {
			found = &i
			break
		}
	}

	if found == nil {
		return ErrTimeSlotNotFound
	}

	if *found == (len(*ds) - 1) {
		// This is the last element remove everything until the index
		*ds = (*ds)[:*found]
		return nil
	}

	// Remove the element at index _found_ and rebuild the slice
	*ds = append((*ds)[:*found], (*ds)[*found+1:]...)
	return nil
}

// At returns a timeslot at a given time
func (ds DailySchedule) At(t Time) *TimeSlot {
	var ret *TimeSlot
	for _, ts := range ds {
		if t.Before(ts.Start) {
			// Return the previously found timeslot
			return ret
		}

		if ts.Stop != nil {
			// Full time slot, reset the ret
			ret = nil
			if t.After(ts.Start) && t.Before(*ts.Stop) {
				// Within the time slot return the time slot with its current
				// index
				return ts
			}
		} else if t.After(ts.Start) {
			// After the beginning of a time range with no stop
			ret = ts
		}
	}

	return ret
}

// FirstAfter returns the first timeslot after a given time
func (ds DailySchedule) FirstAfter(t Time) *TimeSlot {
	for _, ts := range ds {
		if ts.Start.After(t) {
			return ts
		}
	}

	return nil
}

// First returns the first timeslot
func (ds DailySchedule) First() *TimeSlot {
	if len(ds) == 0 {
		return nil
	}

	return ds[0]
}
