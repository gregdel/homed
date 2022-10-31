package xiaomi

import (
	"encoding/json"

	"github.com/gregdel/homed/lib/components"
	"github.com/gregdel/homed/lib/components/common"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	components.Register(components.TypeXiaomiAqara, New)
}

// Climate is a climate sensor
type Climate struct {
	common.Component

	Battery     float64 `json:"battery"`
	Humidity    float64 `json:"humidity"`
	LinkQuality float64 `json:"linkquality"`
	Pressure    float64 `json:"pressure"`
	Temp        float64 `json:"temperature"`
	Voltage     float64 `json:"voltage"`
}

// New returns a new component
func New() components.Component {
	return &Climate{}
}

// Type implements the Component interface
func (c *Climate) Type() components.Type {
	return components.TypeXiaomiAqara
}

// Collectors implements the Component interface
func (c *Climate) Collectors(labels prometheus.Labels) []prometheus.Collector {
	return []prometheus.Collector{
		components.GaugeCollector("temperature", labels,
			func() float64 { return c.Temp },
		),
		components.GaugeCollector("pressure", labels,
			func() float64 { return c.Pressure },
		),
		components.GaugeCollector("humidity", labels,
			func() float64 { return c.Humidity },
		),
	}
}

// Update implements the Component interface
func (c *Climate) Update(value []byte) error {
	return json.Unmarshal(value, c)
}

// Temperature implements the TemperatureGetter interface
func (c *Climate) Temperature() (float64, error) {
	if !c.Device().Online {
		return 0, components.ErrDeviceOffline
	}

	return c.Temp, nil
}
