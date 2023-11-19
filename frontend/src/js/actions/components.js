import { request } from "../request";

import { notificationAdd } from "./notifications";

export const componentsFetch = () =>
  request("COMPONENTS_FETCH", "GET", "/components", null, [
    () => notificationAdd("Updated", "success", 1, "components_fetched"),
  ]);

export const componentUpdate = (id, data) =>
  request("COMPONENT_UPDATE", "PUT", "/components/" + id, data);

export const eventComponentUpdate = (component) => ({
  type: "EVENT_COMPONENT_UPDATE",
  payload: component,
});
