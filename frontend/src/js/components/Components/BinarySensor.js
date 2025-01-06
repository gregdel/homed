import React from "react";
import PropTypes from "prop-types";

import { mdiToggleSwitchOffOutline, mdiToggleSwitch } from "@mdi/js";
import Icon from "@mdi/react";

import { getSwitch } from "./common/switch.js";

export const BinarySensor = ({ id }) => {
  const { on } = getSwitch(id);

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
