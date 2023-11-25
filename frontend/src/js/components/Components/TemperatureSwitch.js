import React from "react";
import { useDispatch, useSelector } from "react-redux";
import PropTypes from "prop-types";

import Icon from "@mdi/react";
import { mdiThermometer, mdiThermometerOff } from "@mdi/js";

import { componentUpdate } from "../../actions/components";

export const TemperatureSwitch = ({ id }) => {
  const dispatch = useDispatch();
  const data = useSelector((state) => state.components.components.get(id));
  if (!data) {
    return null;
  }

  const on = data.values.on;

  const toggle = () => {
    dispatch(componentUpdate(id, on ? "OFF" : "ON"));
  };

  return (
    <div onClick={toggle}>
      <Icon
        path={on ? mdiThermometer : mdiThermometerOff}
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
