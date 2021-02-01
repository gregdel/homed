import { configureAxios, request } from "../request";

export const componentsFetch = () =>
  request("COMPONENTS_FETCH", configureAxios().get("/components"));

export const componentUpdate = (id, data) =>
  request("COMPONENT_UPDATE", configureAxios().put("/components/" + id, data));

export const eventComponentUpdate = (component) => ({
  type: "EVENT_COMPONENT_UPDATE",
  payload: component,
});
