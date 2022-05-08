package components

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

// UpdateFloat64 parses and updates a float64 from an input
func UpdateFloat64(input []byte, output *float64) error {
	v, err := strconv.ParseFloat(string(input), 64)
	if err != nil {
		return err
	}
	*output = v
	return nil
}

// SingleCollector returns a single prometheus collector
func SingleCollector(s Component, labels prometheus.Labels, fn func() float64) []prometheus.Collector {
	return []prometheus.Collector{
		prometheus.NewGaugeFunc(
			prometheus.GaugeOpts{
				Name:        "homed_" + string(s.Type()),
				ConstLabels: labels,
			},
			fn,
		),
	}
}

// SingleBoolCollector returns a single prometheus collector
func SingleBoolCollector(s Component, labels prometheus.Labels, fn func() bool) []prometheus.Collector {
	convert := func() float64 {
		if fn() {
			return 1
		}
		return 0
	}

	return SingleCollector(s, labels, convert)
}
