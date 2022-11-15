import React from "react";
import { useSelector } from "react-redux";
import PropTypes from "prop-types";

import { mdiToggleSwitchOffOutline, mdiToggleSwitch } from "@mdi/js";

import Icon from "@mdi/react";

export const BinarySensor = ({ id }) => {
  const data = useSelector((state) => state.components.components.get(id));
  if (!data) {
    return null;
  }

  return (
    <div style={{ display: "flex", justifyContent: "space-evenly" }}>
      <Icon
        path={
          data.values.on === true ? mdiToggleSwitch : mdiToggleSwitchOffOutline
        }
        size={2}
      />
      <div style={{ fontSize: "2em", fontWeight: 200 }}>
        <span>{data.values.on ? "ON" : "OFF"}</span>
      </div>
    </div>
  );
};

BinarySensor.propTypes = {
  id: PropTypes.string.isRequired,
};
