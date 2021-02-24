import { configureAxios, request } from "../request";

import { notificationAdd } from "./notifications";

export const fetchSchedule = (id) =>
  request(
    "SCHEDULE_FETCH",
    configureAxios().get(`/components/${id}/schedule`),
    null,
    { id }
  );

export const addSchedule = (id, weekday, data) =>
  request(
    "SCHEDULE_ADD",
    configureAxios().post(`/components/${id}/schedule/daily/${weekday}`, data),
    [() => fetchSchedule(id)],
    { id }
  );

export const addScheduleOverride = (id, data) =>
  request(
    "SCHEDULE_ADD_OVERRIDE",
    configureAxios().post(`/components/${id}/schedule/overrides`, data),
    [() => fetchSchedule(id)],
    { id }
  );

export const deleteSchedule = (componentId, weekday, scheduleId) =>
  request(
    "SCHEDULE_DELETE",
    configureAxios().delete(
      `/components/${componentId}/schedule/daily/${weekday}/${scheduleId}`
    ),
    [() => fetchSchedule(componentId)]
  );

export const deleteScheduleOverride = (componentId, overrideId) =>
  request(
    "SCHEDULE_OVERRIDE_DELETE",
    configureAxios().delete(
      `/components/${componentId}/schedule/overrides/${overrideId}`
    ),
    [() => fetchSchedule(componentId)]
  );

export const setScheduleDefault = (id, value) =>
  request(
    "SCHEDULE_SET_DEFAULT",
    configureAxios().post(`/components/${id}/schedule/default`, {
      value: parseInt(value),
    }),
    [
      () => fetchSchedule(id),
      () =>
        notificationAdd(
          "Schedule default updated",
          "success",
          4,
          "schedule_default_update"
        ),
    ],
    { id }
  );
