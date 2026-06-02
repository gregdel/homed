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

type availabilityCollector struct {
	collector prometheus.Collector
	device    *Device
	exempt    bool
}

func (c availabilityCollector) Describe(ch chan<- *prometheus.Desc) {
	c.collector.Describe(ch)
}

func (c availabilityCollector) Collect(ch chan<- prometheus.Metric) {
	if c.exempt || c.device == nil || c.device.MetricsAvailable() {
		c.collector.Collect(ch)
	}
}

func availabilityAwareCollectors(device *Device, componentType Type, collectors []prometheus.Collector) []prometheus.Collector {
	wrapped := make([]prometheus.Collector, 0, len(collectors))
	for _, collector := range collectors {
		wrapped = append(wrapped, availabilityCollector{
			collector: collector,
			device:    device,
			exempt:    componentType == TypeDeviceStatus,
		})
	}

	return wrapped
}
