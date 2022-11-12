package httpd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gregdel/homed/lib/components"
	"github.com/gregdel/homed/lib/schedule"
	"github.com/julienschmidt/httprouter"
)

func (h *httpd) getScheduledComponent(ps httprouter.Params) (components.Scheduled, error) {
	componentID := ps.ByName("id")
	c, err := h.components.Get(componentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get the component: %s", err)
	}

	sc, ok := c.(components.Scheduled)
	if !ok {
		return nil, fmt.Errorf("this component can not be scheduled")
	}

	return sc, nil
}

func (h *httpd) httpGetSchedule(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	sc, err := h.getScheduledComponent(ps)
	if err != nil {
		h.httpError(w, err.Error())
		return
	}

	data := struct {
		Schedule     *schedule.Schedule `json:"schedule"`
		ScheduleName string             `json:"schedule_name"`
	}{
		Schedule:     sc.Schedule(),
		ScheduleName: sc.ScheduleName(),
	}

	h.httpRenderJSON(w, data)
}

func (h *httpd) httpPostSchedule(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ts := schedule.TimeSlot{}
	err := json.NewDecoder(r.Body).Decode(&ts)
	if err != nil {
		h.httpError(w, fmt.Sprintf("failed to decode data: %s", err.Error()))
		return
	}

	weekdayStr := ps.ByName("weekday")
	weekday, err := strconv.Atoi(weekdayStr)
	if err != nil {
		h.httpError(w, fmt.Sprintf("failed to parse weekday: %s", err.Error()))
		return
	}

	if weekday < 0 || weekday > 6 {
		h.httpError(w, "invalid weekday")
		return
	}

	sc, err := h.getScheduledComponent(ps)
	if err != nil {
		h.httpError(w, err.Error())
		return
	}

	schedule := sc.Schedule()
	err = schedule.Add(time.Weekday(weekday), &ts)
	if err != nil {
		h.httpError(w, fmt.Sprintf("failed to add to the schedule: %s", err.Error()))
		return
	}

	err = sc.SaveSchedule(schedule)
	if err != nil {
		h.httpError(w, fmt.Sprintf("failed to save schedule: %s", err.Error()))
		return
	}
}

func (h *httpd) httpPostScheduleDefault(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	d := struct {
		Value float64 `json:"value"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		h.httpError(w, fmt.Sprintf("failed to decode data: %s", err.Error()))
		return
	}

	sc, err := h.getScheduledComponent(ps)
	if err != nil {
		h.httpError(w, err.Error())
		return
	}

	schedule := sc.Schedule()
	schedule.DefaultValue = d.Value
	err = sc.SaveSchedule(schedule)
	if err != nil {
		h.httpError(w, fmt.Sprintf("failed to save schedule: %s", err.Error()))
		return
	}
}

func (h *httpd) httpPostScheduleOverrides(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	override := &schedule.Override{}
	err := json.NewDecoder(r.Body).Decode(override)
	if err != nil {
		h.httpError(w, fmt.Sprintf("failed to decode data: %s", err.Error()))
		return
	}

	sc, err := h.getScheduledComponent(ps)
	if err != nil {
		h.httpError(w, err.Error())
		return
	}

	schedule := sc.Schedule()
	err = schedule.AddOverride(override)
	if err != nil {
		h.httpError(w, fmt.Sprintf("failed to add schedule override: %s", err))
		return
	}

	err = sc.SaveSchedule(schedule)
	if err != nil {
		h.httpError(w, fmt.Sprintf("failed to save schedule: %s", err.Error()))
		return
	}
}

func (h *httpd) httpDeleteSchedule(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	weekdayStr := ps.ByName("weekday")
	weekday, err := strconv.Atoi(weekdayStr)
	if err != nil {
		h.httpError(w, fmt.Sprintf("failed to parse weekday: %s", err.Error()))
		return
	}

	if weekday < 0 || weekday > 6 {
		h.httpError(w, "invalid weekday")
		return
	}

	uuid := ps.ByName("uuid")

	sc, err := h.getScheduledComponent(ps)
	if err != nil {
		h.httpError(w, err.Error())
		return
	}

	schedule := sc.Schedule()
	err = schedule.Delete(time.Weekday(weekday), uuid)
	if err != nil {
		h.httpError(w, fmt.Sprintf("failed to delete the schedule: %s", err.Error()))
		return
	}

	err = sc.SaveSchedule(schedule)
	if err != nil {
		h.httpError(w, fmt.Sprintf("failed to save schedule: %s", err.Error()))
		return
	}
}

func (h *httpd) httpDeleteScheduleOverride(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	sc, err := h.getScheduledComponent(ps)
	if err != nil {
		h.httpError(w, err.Error())
		return
	}

	id := ps.ByName("overrideID")
	schedule := sc.Schedule()
	err = schedule.DeleteOverride(id)
	if err != nil {
		h.httpError(w, fmt.Sprintf("failed to delete the schedule override: %s", err))
		return
	}

	err = sc.SaveSchedule(schedule)
	if err != nil {
		h.httpError(w, fmt.Sprintf("failed to save schedule: %s", err))
		return
	}
}
