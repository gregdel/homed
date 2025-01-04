import React from "react";

import PropTypes from "prop-types";

import { mdiFan } from "@mdi/js";

import { IconToggle } from "./common/IconToggle";

import { useComponents } from "../ComponentsContext";

export const BinaryFan = ({ id }) => {
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
      iconOn={mdiFan}
      iconOff={mdiFan}
      toggle={toggle}
      on={on}
      rotate={on}
    />
  );
};

BinaryFan.propTypes = {
  id: PropTypes.string.isRequired,
};
