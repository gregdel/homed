import { configureAxios, request } from "../request";

export const fetchSchedule = () =>
  request("FETCH_SCHEDULE", configureAxios().get("/schedule"));

export const addSchedule = (data) =>
  request("ADD_SCHEDULE", configureAxios().post("/schedule", data), [
    () => fetchSchedule(),
  ]);

export const deleteSchedule = (weekday, id) =>
  request(
    "DELETE_SCHEDULE",
    configureAxios().delete(`/schedule/${weekday}/${id}`),
    [() => deleteSchedule()]
  );
