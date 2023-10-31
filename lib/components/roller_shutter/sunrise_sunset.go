package rollershutter

import (
	"time"

	"github.com/nathan-osman/go-sunrise"
)

// TwilightType represents the angle of the sun below the horizon
type TwilightType float64

// TwilightType common constants
const (
	TwilightCivil        TwilightType = -6
	TwilightNautical     TwilightType = -12
	TwilightAstronomical TwilightType = -18
)

func (rs *RollerShutter) getSunriseSunset(t time.Time) (time.Time, time.Time) {
	return sunrise.TimeOfElevation(
		rs.Params.Location.Latitude,
		rs.Params.Location.Longitude,
		float64(TwilightCivil),
		t.Year(),
		t.Month(),
		t.Day(),
	)
}
