import React from "react";
import PropTypes from "prop-types";
import { useComponents } from "../ComponentsContext";

import { mdiPower } from "@mdi/js";

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

  return (
    <IconToggle iconOn={mdiPower} iconOff={mdiPower} toggle={toggle} on={on} />
  );
};

Switch.propTypes = {
  id: PropTypes.string.isRequired,
};
