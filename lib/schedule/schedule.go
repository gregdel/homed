package schedule

import (
	"sync"
	"time"
)

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
