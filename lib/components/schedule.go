package components

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/gregdel/homed/lib/schedule"
	"gopkg.in/yaml.v3"
)

const defaultScheduleValue = 14

// Scheduled is an interface to handle the schedule of a component
type Scheduled interface {
	Schedule() *schedule.Schedule
	ScheduleName() string
	SetSchedule(*schedule.Schedule, string, string)
	SaveSchedule() error
}

func schedulePath(path, name string) string {
	return filepath.Join(path, "schedule_"+name+".yaml")
}

func loadSchedule(path string) (*schedule.Schedule, error) {
	s := schedule.New(defaultScheduleValue, false)
	file, err := os.Open(path)
	if err != nil {
		pathError := &os.PathError{}
		if errors.As(err, &pathError) {
			// TODO: should we do this here ? Probably not :)
			return s, nil
		}

		return nil, err
	}
	defer file.Close()

	return s, yaml.NewDecoder(file).Decode(s)
}
