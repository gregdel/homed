package components

import "github.com/prometheus/client_golang/prometheus"

// GaugeCollector returns a prometheus gauge collector
func GaugeCollector(name string, labels prometheus.Labels,
	fn func() float64) prometheus.Collector {
	return prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name:        "homed_" + name,
			ConstLabels: labels,
		}, fn)
}

// CounterCollector returns a prometheus counter collector
func CounterCollector(name string, labels prometheus.Labels,
	fn func() float64) prometheus.Collector {
	return prometheus.NewCounterFunc(
		prometheus.CounterOpts{
			Name:        "homed_" + name,
			ConstLabels: labels,
		}, fn)
}
