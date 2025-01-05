import { combineReducers } from "redux";

// Immer
import { enableMapSet } from "immer";
enableMapSet();

import schedules from "./schedules";

export default combineReducers({
  schedules,
});
