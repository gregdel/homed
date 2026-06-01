package homedtemperature

import "github.com/gregdel/homed/lib/components"

func (h *HomedTemperature) handleBinaryTRV() {
	h.mu.RLock()
	state := h.stateSnapshotLocked()
	target := h.temperatureTargetLocked()
	trvs := make([]components.Switch, 0, len(h.binTRVs))
	for _, trv := range h.binTRVs {
		trvs = append(trvs, trv)
	}
	h.mu.RUnlock()

	if len(trvs) == 0 {
		return
	}

	if !state.On {
		for _, trv := range trvs {
			trv.TurnOff()
		}
		return
	}

	for _, trv := range trvs {
		if state.Current < target {
			trv.TurnOn()
		} else {
			trv.TurnOff()
		}
	}
}
