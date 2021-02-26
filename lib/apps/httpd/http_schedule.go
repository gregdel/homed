package httpd

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gregdel/homed/lib/components"
	"github.com/gregdel/homed/lib/schedule"
	"github.com/julienschmidt/httprouter"
)

var (
	errComponentCannotBeScheduled = errors.New("homed: component cannot be scheduled")
)

func (h *httpd) httpGetSchedule(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	componentID := ps.ByName("id")
	c, err := h.components.Get(componentID)
	if err != nil {
		h.httpError(w, fmt.Sprintf("failed to get the component: %s", err))
		return
	}

	sc, ok := c.(components.Scheduled)
	if !ok {
		h.httpError(w, "this component can not be scheduled")
		return
	}

	h.httpRenderJSON(w, sc.Schedule())
}

func (h *httpd) httpPostSchedule(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	componentID := ps.ByName("id")
	c, err := h.components.Get(componentID)
	if err != nil {
		h.httpError(w, fmt.Sprintf("failed to get the component: %s", err))
		return
	}

	sc, ok := c.(components.Scheduled)
	if !ok {
		h.httpError(w, err.Error())
		return
	}

	ts := schedule.TimeSlot{}
	err = json.NewDecoder(r.Body).Decode(&ts)
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

	schedule := sc.Schedule()
	err = schedule.Add(time.Weekday(weekday), &ts)
	if err != nil {
		h.httpError(w, fmt.Sprintf("failed to add to the schedule: %s", err.Error()))
		return
	}

	err = sc.SaveSchedule(h.components.SchedulePath(c))
	if err != nil {
		h.httpError(w, fmt.Sprintf("failed to save schedule: %s", err.Error()))
		return
	}
}

func (h *httpd) httpPostScheduleDefault(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	componentID := ps.ByName("id")
	c, err := h.components.Get(componentID)
	if err != nil {
		h.httpError(w, fmt.Sprintf("failed to get the component: %s", err))
		return
	}

	sc, ok := c.(components.Scheduled)
	if !ok {
		h.httpError(w, err.Error())
		return
	}

	d := struct {
		Value float64 `json:"value"`
	}{}

	err = json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		h.httpError(w, fmt.Sprintf("failed to decode data: %s", err.Error()))
		return
	}

	schedule := sc.Schedule()
	schedule.DefaultValue = d.Value
	err = sc.SaveSchedule(h.components.SchedulePath(c))
	if err != nil {
		h.httpError(w, fmt.Sprintf("failed to save schedule: %s", err.Error()))
		return
	}
}

func (h *httpd) httpPostScheduleOverrides(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	componentID := ps.ByName("id")
	c, err := h.components.Get(componentID)
	if err != nil {
		h.httpError(w, fmt.Sprintf("failed to get the component: %s", err))
		return
	}

	sc, ok := c.(components.Scheduled)
	if !ok {
		h.httpError(w, err.Error())
		return
	}

	override := &schedule.Override{}
	err = json.NewDecoder(r.Body).Decode(override)
	if err != nil {
		h.httpError(w, fmt.Sprintf("failed to decode data: %s", err.Error()))
		return
	}

	schedule := sc.Schedule()
	err = schedule.AddOverride(override)
	if err != nil {
		h.httpError(w, fmt.Sprintf("failed to add schedule override: %s", err))
		return
	}

	err = sc.SaveSchedule(h.components.SchedulePath(c))
	if err != nil {
		h.httpError(w, fmt.Sprintf("failed to save schedule: %s", err.Error()))
		return
	}
}

func (h *httpd) httpDeleteSchedule(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	componentID := ps.ByName("id")
	c, err := h.components.Get(componentID)
	if err != nil {
		h.httpError(w, fmt.Sprintf("failed to get the component: %s", err.Error()))
		return
	}

	sc, ok := c.(components.Scheduled)
	if !ok {
		h.httpError(w, "this component can not be scheduled")
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

	uuid := ps.ByName("uuid")
	schedule := sc.Schedule()
	err = schedule.Delete(time.Weekday(weekday), uuid)
	if err != nil {
		h.httpError(w, fmt.Sprintf("failed to delete the schedule: %s", err.Error()))
		return
	}

	err = sc.SaveSchedule(h.components.SchedulePath(c))
	if err != nil {
		h.httpError(w, fmt.Sprintf("failed to save schedule: %s", err.Error()))
		return
	}
}

func (h *httpd) httpDeleteScheduleOverride(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	componentID := ps.ByName("id")
	c, err := h.components.Get(componentID)
	if err != nil {
		h.httpError(w, fmt.Sprintf("failed to get the component: %s", err))
		return
	}

	sc, ok := c.(components.Scheduled)
	if !ok {
		h.httpError(w, "this component can not be scheduled")
		return
	}

	id := ps.ByName("overrideID")
	schedule := sc.Schedule()
	err = schedule.DeleteOverride(id)
	if err != nil {
		h.httpError(w, fmt.Sprintf("failed to delete the schedule override: %s", err))
		return
	}

	err = sc.SaveSchedule(h.components.SchedulePath(c))
	if err != nil {
		h.httpError(w, fmt.Sprintf("failed to save schedule: %s", err))
		return
	}
}
