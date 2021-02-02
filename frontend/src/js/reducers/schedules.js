import { produce } from "immer";

const defaultState = {
  schedules: new Map(),
};

export default (state = defaultState, action) =>
  produce(state, (draft) => {
    switch (action.type) {
      case "SCHEDULE_FETCH_FULFILLED":
        draft.schedules.set(action.payload.main.id, action.payload.response);
        break;

      default:
        return draft;
    }
  });
