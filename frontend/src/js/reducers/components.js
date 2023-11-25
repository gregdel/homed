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
          const id = component.values.id;
          const room = component.values.room_name;
          draft.components.set(id, component);

          if (component.type === "homed_temperature") {
            draft.temperatureControl.rooms.set(room, id);
          }

          if (component.type === "homed_temperature_switch") {
            draft.temperatureControl.temperatureSwitch = id;
          }

          if (component.type === "boiler") {
            draft.temperatureControl.boiler = id;
          }
        });

        break;

      case "EVENT_COMPONENT_UPDATE":
        data = action.payload;
        if (!data.values || !data.values.id) {
          return draft;
        }

        // TODO: we only update the values of the sensor to avoid losing the
        // romm / device information, we should find a better way to handle
        // this
        var current = draft.components.get(data.values.id);
        current.values = data.values;

        draft.components.set(data.values.id, current);
        break;
      default:
        return draft;
    }
  });
