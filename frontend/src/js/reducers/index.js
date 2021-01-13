import { combineReducers } from "redux";

// Immer
import { enableMapSet } from "immer";
enableMapSet();

import stuff from "./stuff";

export default combineReducers({
  stuff,
});
