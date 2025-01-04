import { combineReducers } from "redux";

// Immer
import { enableMapSet } from "immer";
enableMapSet();

import schedules from "./schedules";
import notifications from "./notifications";

export default combineReducers({
  schedules,
  notifications,
});
