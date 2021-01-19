package schedule

import (
	"sync"
	"time"
)

// Represents the max number of days to search for
const maxDaysSearch = 7

// Schedule holds the schedule
type Schedule struct {
	mu   sync.Mutex
	Days map[time.Weekday]*DailySchedule `json:"days"`
}

// New returns a new schedule
func New() *Schedule {
	schedule := &Schedule{
		Days: map[time.Weekday]*DailySchedule{},
	}

	for i := 0; i < 7; i++ {
		ds := NewDailySchedule()
		schedule.Days[time.Weekday(i)] = &ds
	}

	return schedule
}

// At returns the timestlot at a given time
func (s *Schedule) At(day time.Weekday, time Time) *TimeSlot {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Days[day].At(time)
}

// Now returns the current timeslot
func (s *Schedule) Now() *TimeSlot {
	now := now()
	return s.At(
		now.Weekday(),
		NewTime(now.Hour(), now.Minute(), now.Second()),
	)
}

// Next returns the next timeslot
func (s *Schedule) Next() (*TimeSlot, time.Weekday) {
	now := now()

	weekday := now.Weekday()
	currentDay := true

	for i := 0; i < maxDaysSearch; i++ {
		d := time.Weekday((int(weekday) + i) % 7)
		ds := s.Days[d]

		var ts *TimeSlot
		if currentDay {
			ts = ds.FirstAfter(NewTime(now.Hour(), now.Minute(), now.Second()))
		} else {
			ts = ds.First()
		}

		if ts != nil {
			return ts, d
		}

		currentDay = false
	}

	return nil, 0
}

// NextTime returns the timeslot and time of the next timeslot
func (s *Schedule) NextTime() (*TimeSlot, *time.Time) {
	ts, wd := s.Next()
	if ts == nil {
		return nil, nil
	}

	now := now()
	day := wd - now.Weekday()
	if day < 0 {
		day += 7
	}

	days := time.Duration(int(day)*24) * time.Hour
	hours := time.Duration(ts.Start.Hour-now.Hour()) * time.Hour
	minutes := time.Duration(ts.Start.Minute-now.Minute()) * time.Minute
	seconds := time.Duration(ts.Start.Second-now.Second()) * time.Second
	t := now.Add(days + hours + minutes + seconds)

	return ts, &t
}

// Add adds a timeslot to a schedule
func (s *Schedule) Add(day time.Weekday, slot *TimeSlot) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	ds := s.Days[day]
	return ds.Add(slot)
}

// Delete deletes a timeslot on a given day
func (s *Schedule) Delete(day time.Weekday, uuid string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	ds := s.Days[day]
	return ds.Delete(uuid)
}
