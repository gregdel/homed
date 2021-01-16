import { createHashHistory } from "history";
import { createStore, applyMiddleware, compose } from "redux";
import thunk from "redux-thunk";

export const history = createHashHistory();

import rootReducer from "./reducers/index";

const middlewares = [thunk];

// Only use in development mode (set in webpack)
// if (process.env.NODE_ENV === "development") {
//   const { logger } = require("redux-logger");
//   middlewares.push(logger);
// }

// Export the store
const store = createStore(
  rootReducer,
  compose(applyMiddleware(...middlewares))
);

export default store;
