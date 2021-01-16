import { produce } from "immer";

const defaultState = {
  schedule: {},
};

export default (state = defaultState, action) =>
  produce(state, (draft) => {
    var data;
    switch (action.type) {
      case "FETCH_SCHEDULE_FULFILLED":
        data = action.payload.response;
        draft.schedule = data;
        break;
      default:
        return draft;
    }
  });
