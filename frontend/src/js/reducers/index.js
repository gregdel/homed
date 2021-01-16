import { combineReducers } from "redux";

// Immer
import { enableMapSet } from "immer";
enableMapSet();

import stuff from "./stuff";
import temperatureSchedule from "./temperatureSchedule";

export default combineReducers({
  stuff,
  temperatureSchedule,
});
