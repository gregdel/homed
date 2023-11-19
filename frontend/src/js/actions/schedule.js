import { request } from "../request";

import { notificationAdd } from "./notifications";

export const fetchSchedule = (id) =>
  request("SCHEDULE_FETCH", "GET", `/components/${id}/schedule`, null, null, {
    id,
  });

export const addSchedule = (id, weekday, data) =>
  request(
    "SCHEDULE_ADD",
    "POST",
    `/components/${id}/schedule/daily/${weekday}`,
    data,
    [() => fetchSchedule(id)],
    { id }
  );

export const addScheduleOverride = (id, data) =>
  request(
    "SCHEDULE_ADD_OVERRIDE",
    "POST",
    `/components/${id}/schedule/overrides`,
    data,
    [() => fetchSchedule(id)],
    { id }
  );

export const deleteSchedule = (componentId, weekday, scheduleId) =>
  request(
    "SCHEDULE_DELETE",
    "DELETE",
    `/components/${componentId}/schedule/daily/${weekday}/${scheduleId}`,
    null,
    [() => fetchSchedule(componentId)]
  );

export const deleteScheduleOverride = (componentId, overrideId) =>
  request(
    "SCHEDULE_OVERRIDE_DELETE",
    "DELETE",
    `/components/${componentId}/schedule/overrides/${overrideId}`,
    null,
    [() => fetchSchedule(componentId)]
  );

export const setScheduleDefault = (id, value) =>
  request(
    "SCHEDULE_SET_DEFAULT",
    "POST",
    `/components/${id}/schedule/default`,
    { value: parseInt(value) },
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
