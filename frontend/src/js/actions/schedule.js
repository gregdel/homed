import { configureAxios, request } from "../request";

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
    configureAxios().post(`/components/${id}/schedule/${weekday}`, data),
    [() => fetchSchedule(id)],
    { id }
  );

export const deleteSchedule = (componentId, weekday, scheduleId) =>
  request(
    "SCHEDULE_DELETE",
    configureAxios().delete(
      `/components/${componentId}/schedule/${weekday}/${scheduleId}`
    ),
    [() => fetchSchedule(componentId)]
  );
