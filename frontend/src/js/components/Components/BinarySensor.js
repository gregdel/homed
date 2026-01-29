import React from "react";
import PropTypes from "prop-types";

import { Icon } from "../ui/Icon";
import { getSwitch } from "./common/switch.js";

export const BinarySensor = ({ id }) => {
  const { on } = getSwitch(id);

  return (
    <div style={{ display: "flex", justifyContent: "space-evenly" }}>
      <Icon name={on ? "toggleSwitch" : "toggleSwitchOffOutline"} size={2} />
      <div style={{ fontSize: "2em", fontWeight: 200 }}>
        <span>{on ? "ON" : "OFF"}</span>
      </div>
    </div>
  );
};

BinarySensor.propTypes = {
  id: PropTypes.string.isRequired,
};
