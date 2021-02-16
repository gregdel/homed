package common

import (
	"os"
	"time"

	"github.com/gregdel/homed/lib/schedule"
	"gopkg.in/yaml.v2"
)

// Const for now
const defaultValue = 14

// ScheduledComponent represents a scheduled component
type ScheduledComponent struct {
	Component

	schedule *schedule.Schedule
}

// NewScheduledComponent returns a new scheduled component
func NewScheduledComponent() *ScheduledComponent {
	return &ScheduledComponent{
		schedule: schedule.New(defaultValue),
	}
}

// ScheduledNextTime implements the Scheduled interface
func (c *ScheduledComponent) ScheduledNextTime() *time.Time {
	if c.schedule == nil {
		return nil
	}

	_, nt := c.schedule.NextTime()
	return nt
}

// ScheduledDefault implements the Scheduled interface
func (c *ScheduledComponent) ScheduledDefault() float64 {
	return c.schedule.DefaultValue
}

// CurrentSchedule implements the Scheduled interface
func (c *ScheduledComponent) CurrentSchedule() float64 {
	if c.schedule == nil {
		return -1
	}

	return c.schedule.Value()
}

// Schedule implements the Scheduled interface
func (c *ScheduledComponent) Schedule() *schedule.Schedule {
	return c.schedule
}

// SetSchedule implements the Scheduled interface
func (c *ScheduledComponent) SetSchedule(s *schedule.Schedule) {
	if s == nil {
		return
	}

	c.schedule = s
}

// ScheduleAdd implements the Scheduled interface
func (c *ScheduledComponent) ScheduleAdd(wd time.Weekday, ts *schedule.TimeSlot) error {
	return c.schedule.Add(wd, ts)
}

// ScheduleDelete implements the Scheduled interface
func (c *ScheduledComponent) ScheduleDelete(wd time.Weekday, id string) error {
	return c.schedule.Delete(wd, id)
}

// LoadSchedule implements the Scheduled interface
func (c *ScheduledComponent) LoadSchedule(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return yaml.NewDecoder(file).Decode(c.schedule)
}

// SaveSchedule implements the Scheduled interface
func (c *ScheduledComponent) SaveSchedule(path string) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	return yaml.NewEncoder(file).Encode(c.schedule)
}
