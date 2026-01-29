import React from "react";
import PropTypes from "prop-types";
import { useComponents } from "../ComponentsContext";

import { Icon } from "../ui/Icon";

export const TemperatureSwitch = ({ id }) => {
  const { getComponentById, updateComponent } = useComponents();

  const component = getComponentById(id);
  if (component === undefined) {
    return null;
  }

  const { on } = component.values;

  const toggle = () => {
    updateComponent(id, on ? "OFF" : "ON");
  };

  return (
    <div onClick={toggle}>
      <Icon
        name={on ? "thermometer" : "thermometerOff"}
        size={3}
        style={{
          cursor: "pointer",
          color: on ? "#000000" : "#00000040",
          transition: "color 0.3s ease-out 0s",
        }}
      />
    </div>
  );
};

TemperatureSwitch.propTypes = {
  id: PropTypes.string.isRequired,
};
