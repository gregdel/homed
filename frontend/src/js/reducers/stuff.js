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
      case "FETCH_STUFF_FULFILLED":
        data = action.payload.response;

        data.rooms.map((room) => {
          var devices = room.devices;
          devices.map((device) => {
            var components = device.components;
            components.map((component) => {
              const uuid = component.values.uuid;

              draft.components.set(uuid, {
                room: room.name,
                device: device.name,
                ...component,
              });

              // Keep the temperature controled rooms in a different map
              if (component.type === "homed_temperature") {
                draft.temperatureControl.rooms.set(room.name, uuid);
              }

              // Keep the temperature controled rooms in a different map
              if (component.type === "boiler") {
                draft.temperatureControl.boiler = uuid;
              }
            });
          });
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
