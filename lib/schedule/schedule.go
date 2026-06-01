package schedule

import (
	"sync"
	"time"
)

// Represents the max number of days to search for
const maxDaysSearch = 7

// Schedule holds the schedule
type Schedule struct {
	mu sync.Mutex

	Days         map[time.Weekday]*DailySchedule `json:"days" yaml:"days"`
	DefaultValue float64                         `json:"default_value" yaml:"default_value"`
	DefaultOn    bool                            `json:"default_on" yaml:"default_on"`
	Overrides    Overrides                       `json:"overrides" yaml:"overrides"`
}

// New returns a new schedule
func New(dv float64, do bool) *Schedule {
	schedule := &Schedule{
		Days:         map[time.Weekday]*DailySchedule{},
		DefaultValue: dv,
		DefaultOn:    do,
		Overrides:    NewOverrides(),
	}

	for i := range 7 {
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

// nextTimeslot returns the next timeslot
func (s *Schedule) nextTimeslot() (*TimeSlot, time.Weekday) {
	now := now()

	weekday := now.Weekday()
	currentDay := true

	for i := range maxDaysSearch {
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

func (s *Schedule) newTime(t Time, wd time.Weekday) *time.Time {
	now := now()
	day := wd - now.Weekday()
	if day < 0 {
		day += 7
	}

	days := time.Duration(int(day)*24) * time.Hour
	hours := time.Duration(t.Hour-now.Hour()) * time.Hour
	minutes := time.Duration(t.Minute-now.Minute()) * time.Minute
	seconds := time.Duration(t.Second-now.Second()) * time.Second
	scheduledTime := now.Add(days + hours + minutes + seconds)
	return &scheduledTime

}

func (s *Schedule) nextTimeFromTimeSlot() *time.Time {
	ts, wd := s.nextTimeslot()
	if ts == nil {
		return nil
	}

	return s.newTime(ts.Start, wd)
}

// NextChange returns the next time of schedule change
func (s *Schedule) NextChange() *time.Time {
	currentSchedule := s.Now()
	if currentSchedule != nil && currentSchedule.Stop != nil {
		n := now()
		weekday := n.Weekday()
		return s.newTime(*currentSchedule.Stop, weekday)
	}

	return s.nextTimeFromTimeSlot()
}

// Add adds a timeslot to a schedule
func (s *Schedule) Add(day time.Weekday, slot *TimeSlot) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	ds := s.Days[day]
	return ds.Add(slot)
}

// Update updates the schedule by replacing a timeslot by another
func (s *Schedule) Update(day time.Weekday, slot *TimeSlot, id string) error {
	if err := s.Delete(day, id); err != nil {
		return err
	}

	return s.Add(day, slot)
}

// AddOverride adds a schedule override
func (s *Schedule) AddOverride(override *Override) error {
	return s.Overrides.Add(override)
}

// Delete deletes a timeslot on a given day
func (s *Schedule) Delete(day time.Weekday, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	ds := s.Days[day]
	return ds.Delete(id)
}

// DeleteOverride deletes a scheduled override
func (s *Schedule) DeleteOverride(id string) error {
	return s.Overrides.Delete(id)
}

// Values returns the scheduled values at the current time
func (s *Schedule) Values() (float64, bool) {
	if o := s.Overrides.Now(); o != nil {
		return o.Value, o.On
	}

	if ts := s.Now(); ts != nil {
		return ts.Value, ts.On
	}

	return s.DefaultValue, s.DefaultOn
}
