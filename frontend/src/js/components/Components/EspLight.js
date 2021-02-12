import React from "react";
import { useDispatch, useSelector } from "react-redux";

import PropTypes from "prop-types";

import { mdiLightbulbOnOutline, mdiLightbulbOn } from "@mdi/js";

import { IconToggle } from "./common/IconToggle";

import { componentUpdate } from "../../actions/components";

export const ESPLight = ({ id }) => {
  const dispatch = useDispatch();
  const data = useSelector((state) => state.components.components.get(id));
  if (!data) {
    return null;
  }

  const on = data.values.on;

  const toggle = () => {
    const data = { state: on ? "OFF" : "ON" };
    dispatch(componentUpdate(id, data));
  };

  return (
    <IconToggle
      iconOn={mdiLightbulbOn}
      iconOff={mdiLightbulbOnOutline}
      toggle={toggle}
      on={on}
    />
  );
};

ESPLight.propTypes = {
  id: PropTypes.string.isRequired,
};
