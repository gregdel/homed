package homedtemperature

import (
	"log/slog"

	"github.com/gregdel/homed/lib/components"
)

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
			if err := trv.TurnOff(); err != nil {
				h.log.Warn("failed to turn off binary TRV", slog.Any("error", err))
			}
		}
		return
	}

	for _, trv := range trvs {
		if state.Current < target {
			if err := trv.TurnOn(); err != nil {
				h.log.Warn("failed to turn on binary TRV", slog.Any("error", err))
			}
		} else {
			if err := trv.TurnOff(); err != nil {
				h.log.Warn("failed to turn off binary TRV", slog.Any("error", err))
			}
		}
	}
}
