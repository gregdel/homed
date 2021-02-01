import { produce } from "immer";

const defaultState = {
  components: new Map(),
  temperatureControl: {
    rooms: new Map(),
    boiler: "",
  },
};

export default (state = defaultState, action) =>
  produce(state, (draft) => {
    var data;
    switch (action.type) {
      case "COMPONENTS_FETCH_FULFILLED":
        data = action.payload.response;

        data.map((component) => {
          const uuid = component.values.uuid;
          const room = component.values.room_name;
          draft.components.set(uuid, component);

          // Keep the temperature controled rooms in a different map
          if (component.type === "homed_temperature") {
            draft.temperatureControl.rooms.set(room, uuid);
          }

          // Keep the temperature controled rooms in a different map
          if (component.type === "boiler") {
            draft.temperatureControl.boiler = uuid;
          }
        });

        break;

      case "EVENT_COMPONENT_UPDATE":
        data = action.payload;
        if (!data.values || !data.values.uuid) {
          return draft;
        }

        // TODO: we only update the values of the sensor to avoid losing the
        // romm / device information, we should find a better way to handle
        // this
        var current = draft.components.get(data.values.uuid);
        current.values = data.values;

        draft.components.set(data.values.uuid, current);
        break;
      default:
        return draft;
    }
  });
