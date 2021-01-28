package scheduled

import (
	"time"

	"github.com/gregdel/homed/lib/schedule"
)

// ScheduledComponent represents a scheduled component
type ScheduledComponent struct {
	schedule *schedule.Schedule
}

// New reuturns a new scheduled component
func New() *ScheduledComponent {
	return &ScheduledComponent{schedule: schedule.New()}
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
	return 14
}

// CurrentSchedule implements the Scheduled interface
func (c *ScheduledComponent) CurrentSchedule() float64 {
	if c.schedule == nil {
		return -1
	}

	ts := c.schedule.Now()
	if ts != nil {
		return ts.Value
	}

	return c.ScheduledDefault()
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
