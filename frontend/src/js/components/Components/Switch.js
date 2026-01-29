import React from "react";
import PropTypes from "prop-types";
import { useComponents } from "../ComponentsContext";

import { IconToggle } from "./common/IconToggle";

export const Switch = ({ id }) => {
  const { getComponentById, updateComponent } = useComponents();

  const component = getComponentById(id);
  if (component === undefined) {
    return null;
  }

  const { on } = component.values;

  const toggle = () => {
    updateComponent(id, on ? "OFF" : "ON");
  };

  return <IconToggle iconOn="power" iconOff="power" toggle={toggle} on={on} />;
};

Switch.propTypes = {
  id: PropTypes.string.isRequired,
};
