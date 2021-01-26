package rtl433

import (
	"encoding/json"

	"github.com/gregdel/homed/lib/components"
	base "github.com/gregdel/homed/lib/components/base_component"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	components.Register(components.TypeRTL433, New)
}

// RTL433 is a component that handles rtl_433 signals
type RTL433 struct {
	base.Component

	BoilerState bool `json:"boiler_state"`
}

// New returns a new humidity component
func New() components.Component {
	return &RTL433{}
}

// Type implements the Component interface
func (r *RTL433) Type() components.Type {
	return components.TypeRTL433
}

// Collectors implements the Component interface
func (r *RTL433) Collectors(labels prometheus.Labels) []prometheus.Collector {
	return []prometheus.Collector{
		prometheus.NewGaugeFunc(
			prometheus.GaugeOpts{Name: "homed_boiler_state"},
			func() float64 {
				if r.BoilerState {
					return 1
				}
				return 0
			},
		),
	}
}

// Update implements the Component interface
func (r *RTL433) Update(value []byte) error {
	data := struct {
		Model string `json:"model"`
		Rows  []struct {
			Len  int    `json:"len"`
			Data string `json:"data"`
		} `json:"rows"`
	}{}

	if err := json.Unmarshal(value, &data); err != nil {
		return nil
	}

	if data.Model != "boiler" {
		return nil
	}

	if len(data.Rows) != 3 {
		return nil
	}

	code := data.Rows[2].Data
	switch code {
	case "ff67b7efb7fbb8":
		r.BoilerState = true
	case "ff67b7ffdbfddf":
		r.BoilerState = true
	case "ff67b7ffdbfedf":
		r.BoilerState = false
	case "ff67b7efb7fdb8":
		r.BoilerState = false
	}

	return nil
}
