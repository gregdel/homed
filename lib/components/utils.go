package components

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

func updateFloat64(input []byte, output *float64) error {
	v, err := strconv.ParseFloat(string(input), 64)
	if err != nil {
		return err
	}
	*output = v
	return nil
}

func singleCollector(s Component, labels prometheus.Labels, fn func() float64) []prometheus.Collector {
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
