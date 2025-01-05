import { request } from "../request";

export const fetchSchedule = (cid) =>
  request("SCHEDULE_FETCH", "GET", `/components/${cid}/schedule`, null, null, {
    id: cid,
  });

export const addScheduleTimeSlot = (cid, weekday, data) =>
  request(
    "SCHEDULE_ADD_TIMESLOT",
    "POST",
    `/components/${cid}/schedule/daily/${weekday}`,
    data,
    [() => fetchSchedule(cid)]
  );

export const updateScheduleTimeSlot = (cid, weekday, data, id) =>
  request(
    "SCHEDULE_UPDATE_TIMESLOT",
    "PUT",
    `/components/${cid}/schedule/daily/${weekday}/${id}`,
    data,
    [() => fetchSchedule(cid)]
  );

export const deleteScheduleTimeSlot = (cid, weekday, id) =>
  request(
    "SCHEDULE_DELETE_TIMESLOT",
    "DELETE",
    `/components/${cid}/schedule/daily/${weekday}/${id}`,
    null,
    [() => fetchSchedule(cid)]
  );

export const addScheduleOverride = (cid, data) =>
  request(
    "SCHEDULE_ADD_OVERRIDE",
    "POST",
    `/components/${cid}/schedule/overrides`,
    data,
    [() => fetchSchedule(cid)]
  );

export const deleteScheduleOverride = (cid, oid) =>
  request(
    "SCHEDULE_DELETE_OVERRIDE",
    "DELETE",
    `/components/${cid}/schedule/overrides/${oid}`,
    null,
    [() => fetchSchedule(cid)]
  );

export const setScheduleDefault = (cid, value) =>
  request(
    "SCHEDULE_SET_DEFAULT",
    "POST",
    `/components/${cid}/schedule/default`,
    { value: parseInt(value) },
    [
      () => fetchSchedule(cid),
      // () =>
      //   notificationAdd(
      //     "Schedule default updated",
      //     "success",
      //     4,
      //     "schedule_default_update"
      //   ),
    ]
  );
