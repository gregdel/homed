package homedtemperature

import "go.uber.org/zap"

func (h *HomedTemperature) handleBinaryTRV() {
	if len(h.binTRVs) == 0 || !h.On.Load() {
		return
	}

	current := h.Current.Load()
	target, err := h.TemperatureTarget()
	if err != nil {
		h.log.Warn("failed to get temperature target", zap.Error(err))
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
