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
		return nil, fmt.Errorf("failed to get the component: %w", err)
	}

	sc, ok := c.(components.Scheduled)
	if !ok {
		return nil, fmt.Errorf("this component can not be scheduled")
	}

	return sc, nil
}

func (h *httpd) saveSchedule(w http.ResponseWriter, sc components.Scheduled) {
	if err := sc.SaveSchedule(); err != nil {
		h.httpError(w, fmt.Sprintf("failed to save schedule: %s", err))
		return
	}

	h.httpRenderJSON(w, nil)
}

func (h *httpd) timeSlot(r *http.Request) (*schedule.TimeSlot, error) {
	ts := schedule.TimeSlot{}
	if err := json.NewDecoder(r.Body).Decode(&ts); err != nil {
		return nil, err
	}

	return &ts, nil
}

func (h *httpd) getWeekday(ps httprouter.Params) (time.Weekday, error) {
	weekdayStr := ps.ByName("weekday")
	weekday, err := strconv.Atoi(weekdayStr)
	if err != nil {
		return 0, fmt.Errorf("failed to parse weekday: %w", err)
	}

	if weekday < 0 || weekday > 6 {
		return 0, fmt.Errorf("invalid weekday: %s", weekdayStr)
	}

	return time.Weekday(weekday), nil
}

func (h *httpd) getSchedule(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

func (h *httpd) addScheduleTimeSlot(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ts, err := h.timeSlot(r)
	if err != nil {
		h.httpError(w, fmt.Sprintf("failed to decode timeslot: %s", err.Error()))
		return
	}

	weekday, err := h.getWeekday(ps)
	if err != nil {
		h.httpError(w, err.Error())
		return
	}

	sc, err := h.getScheduledComponent(ps)
	if err != nil {
		h.httpError(w, err.Error())
		return
	}

	err = sc.Schedule().Add(weekday, ts)
	if err != nil {
		h.httpError(w, fmt.Sprintf("failed to add to the schedule: %s", err.Error()))
		return
	}

	h.saveSchedule(w, sc)
}

func (h *httpd) updateScheduleDefault(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	sc.Schedule().DefaultValue = d.Value
	h.saveSchedule(w, sc)
}

func (h *httpd) getScheduleOverrides(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	err = sc.Schedule().AddOverride(override)
	if err != nil {
		h.httpError(w, fmt.Sprintf("failed to add schedule override: %s", err))
		return
	}

	h.saveSchedule(w, sc)
}

func (h *httpd) updateScheduleTimeSlot(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	tsID := ps.ByName("tsID")

	weekday, err := h.getWeekday(ps)
	if err != nil {
		h.httpError(w, err.Error())
		return
	}

	sc, err := h.getScheduledComponent(ps)
	if err != nil {
		h.httpError(w, err.Error())
		return
	}

	ts, err := h.timeSlot(r)
	if err != nil {
		h.httpError(w, fmt.Sprintf("failed to decode timeslot: %s", err.Error()))
		return
	}

	err = sc.Schedule().Update(weekday, ts, tsID)
	if err != nil {
		h.httpError(w, fmt.Sprintf("failed to delete the schedule: %s", err.Error()))
		return
	}

	h.saveSchedule(w, sc)
}

func (h *httpd) deleteScheduleTimeSlot(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	weekday, err := h.getWeekday(ps)
	if err != nil {
		h.httpError(w, err.Error())
		return
	}

	tsID := ps.ByName("tsID")

	sc, err := h.getScheduledComponent(ps)
	if err != nil {
		h.httpError(w, err.Error())
		return
	}

	err = sc.Schedule().Delete(weekday, tsID)
	if err != nil {
		h.httpError(w, fmt.Sprintf("failed to delete the schedule: %s", err.Error()))
		return
	}

	h.saveSchedule(w, sc)
}

func (h *httpd) deleteScheduleOverride(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	sc, err := h.getScheduledComponent(ps)
	if err != nil {
		h.httpError(w, err.Error())
		return
	}

	id := ps.ByName("overrideID")
	err = sc.Schedule().DeleteOverride(id)
	if err != nil {
		h.httpError(w, fmt.Sprintf("failed to delete the schedule override: %s", err))
		return
	}

	h.saveSchedule(w, sc)
}
