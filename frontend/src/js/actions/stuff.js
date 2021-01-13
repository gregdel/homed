import { configureAxios, request } from "../request";

export const fetchStuff = () =>
  request("FETCH_STUFF", configureAxios().get("/data"));
