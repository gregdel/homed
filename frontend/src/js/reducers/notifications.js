import { produce } from "immer";

const defaultState = new Map();

export default (state = defaultState, action) =>
  produce(state, (draft) => {
    switch (action.type) {
      case "NOTIFICATION_ADD":
        draft.set(action.payload.id, action.payload);
        break;

      case "NOTIFICATION_REMOVE":
        draft.delete(action.payload.id);
        break;

      default:
        return draft;
    }
  });
