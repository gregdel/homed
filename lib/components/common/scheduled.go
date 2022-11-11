package common

import (
	"os"

	"github.com/gregdel/homed/lib/schedule"
	"gopkg.in/yaml.v2"
)

// ScheduledComponent represents a scheduled component
type ScheduledComponent struct {
	Component

	scheduleName string
	schedulePath string
	schedule     *schedule.Schedule
}

// NewScheduledComponent returns a new scheduled component
func NewScheduledComponent() *ScheduledComponent {
	return &ScheduledComponent{}
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
func (c *ScheduledComponent) SetSchedule(s *schedule.Schedule, path, name string) {
	c.schedule = s
	c.scheduleName = name
	c.schedulePath = path
}

// SaveSchedule implements the Scheduled interface
func (c *ScheduledComponent) SaveSchedule(s *schedule.Schedule) error {
	file, err := os.OpenFile(c.schedulePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	return yaml.NewEncoder(file).Encode(s)
}
