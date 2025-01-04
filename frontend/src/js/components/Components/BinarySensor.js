import React from "react";
import PropTypes from "prop-types";
import { useComponents } from "../ComponentsContext";

import { mdiToggleSwitchOffOutline, mdiToggleSwitch } from "@mdi/js";
import Icon from "@mdi/react";

export const BinarySensor = ({ id }) => {
  const { getComponentById } = useComponents();

  const component = getComponentById(id);
  if (component === undefined) {
    return null;
  }

  const { on } = component.values;

  return (
    <div style={{ display: "flex", justifyContent: "space-evenly" }}>
      <Icon path={on ? mdiToggleSwitch : mdiToggleSwitchOffOutline} size={2} />
      <div style={{ fontSize: "2em", fontWeight: 200 }}>
        <span>{on ? "ON" : "OFF"}</span>
      </div>
    </div>
  );
};

BinarySensor.propTypes = {
  id: PropTypes.string.isRequired,
};
