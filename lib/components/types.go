package components

import (
	"time"

	"github.com/gregdel/homed/lib/schedule"
)

// Switch is an interface to reprensents a switch
type Switch interface {
	IsOn() bool
	Toggle() error
	SetOn() error
	SetOff() error
	Set(state bool) error
}

// Publisher is an interface to publish the component state
type Publisher interface {
	PublishState() error
}

// Scheduled is an interface to handle the schedule of a component
type Scheduled interface {
	ScheduledNextTime() *time.Time
	ScheduledDefault() float64
	CurrentSchedule() float64
	Schedule() *schedule.Schedule
	SetSchedule(*schedule.Schedule)
	ScheduleAdd(time.Weekday, *schedule.TimeSlot) error
	ScheduleDelete(time.Weekday, string) error
}
