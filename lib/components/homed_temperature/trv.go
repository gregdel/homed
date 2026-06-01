package homedtemperature

import "log/slog"

func (h *HomedTemperature) handleBinaryTRV() {
	if len(h.binTRVs) == 0 {
		return
	}

	if !h.IsOn() {
		for _, trv := range h.binTRVs {
			trv.TurnOff()
		}
		return
	}

	current := h.Current.Load()
	target, err := h.TemperatureTarget()
	if err != nil {
		h.log.Warn("failed to get temperature target", slog.Any("error", err))
		return
	}

	for _, trv := range h.binTRVs {
		if current < target {
			trv.TurnOn()
		} else {
			trv.TurnOff()
		}
	}
}
