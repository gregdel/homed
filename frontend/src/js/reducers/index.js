import { combineReducers } from "redux";

// Immer
import { enableMapSet } from "immer";
enableMapSet();

import components from "./components";
import schedules from "./schedules";
import notifications from "./notifications";

export default combineReducers({
  components,
  schedules,
  notifications,
});
