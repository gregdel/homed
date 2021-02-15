import { configureAxios, request } from "../request";

import { notificationAdd } from "./notifications";

export const componentsFetch = () =>
  request("COMPONENTS_FETCH", configureAxios().get("/components"), [
    () => notificationAdd("Updated", "success", 1, "components_fetched"),
  ]);

export const componentUpdate = (id, data) =>
  request("COMPONENT_UPDATE", configureAxios().put("/components/" + id, data));

export const eventComponentUpdate = (component) => ({
  type: "EVENT_COMPONENT_UPDATE",
  payload: component,
});
