package rollershutter

import (
	"time"

	"github.com/kelvins/sunrisesunset"
	"go.uber.org/zap"
)

func (rs *RollerShutter) getDawnDusk(t time.Time) (time.Time, time.Time) {
	sunrise, sunset := rs.getSunriseSunset(t)

	// Let's say the difference between dawn -> sunrise and sunset -> dusk is
	// around 40 minutes.
	d := 40 * time.Minute

	return sunrise.Add(-1 * d), sunset.Add(d)
}

func (rs *RollerShutter) getSunriseSunset(t time.Time) (time.Time, time.Time) {
	_, utcOffset := t.Zone()
	p := sunrisesunset.Parameters{
		Latitude:  rs.Params.Location.Latitude,
		Longitude: rs.Params.Location.Longitude,
		UtcOffset: float64(utcOffset / 3600),
		Date:      t,
	}

	sunrise, sunset, err := p.GetSunriseSunset()
	if err != nil {
		rs.logger.Error(
			"failed to get sunrise / sunset times, falling back",
			zap.Error(err),
		)

		// Fallback to 8:00 (sunrise) / 19:00 (sunset)
		sunrise = time.Date(t.Year(), t.Month(), t.Day(), 8, 00, 0, 0, t.Location())
		sunset = time.Date(t.Year(), t.Month(), t.Day(), 19, 0, 0, 0, t.Location())
	}

	return sunrise, sunset
}
