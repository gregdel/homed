import { combineReducers } from "redux";

// Immer
import { enableMapSet } from "immer";
enableMapSet();

import components from "./components";
import schedules from "./schedules";
import temperatureSchedule from "./temperatureSchedule";

export default combineReducers({
  components,
  schedules,
  temperatureSchedule,
});
