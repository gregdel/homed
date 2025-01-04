import React from "react";
import PropTypes from "prop-types";
import { useComponents } from "../ComponentsContext";

import { mdiLightbulbOnOutline, mdiLightbulbOn } from "@mdi/js";

import { IconToggle } from "./common/IconToggle";

export const BinaryLight = ({ id }) => {
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
      iconOn={mdiLightbulbOn}
      iconOff={mdiLightbulbOnOutline}
      toggle={toggle}
      on={on}
    />
  );
};

BinaryLight.propTypes = {
  id: PropTypes.string.isRequired,
};
