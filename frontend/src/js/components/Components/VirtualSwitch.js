import React from "react";
import PropTypes from "prop-types";
import { useComponents } from "../ComponentsContext";

import { mdiLightSwitch, mdiLightSwitchOff } from "@mdi/js";

import { IconToggle } from "./common/IconToggle";

export const VirtualSwitch = ({ id }) => {
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
