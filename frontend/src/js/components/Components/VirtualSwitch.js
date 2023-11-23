import React from "react";
import { useDispatch, useSelector } from "react-redux";
import PropTypes from "prop-types";

import { mdiLightSwitch, mdiLightSwitchOff } from "@mdi/js";

import { IconToggle } from "./common/IconToggle";

import { componentUpdate } from "../../actions/components";

export const VirtualSwitch = ({ id }) => {
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
    <IconToggle
      iconOn={mdiLightSwitch}
      iconOff={mdiLightSwitchOff}
      toggle={toggle}
      on={on}
    />
  );
};

VirtualSwitch.propTypes = {
  id: PropTypes.string.isRequired,
};
