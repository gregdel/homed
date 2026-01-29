import React from "react";
import PropTypes from "prop-types";

import { IconToggle } from "./common/IconToggle";
import { getSwitch } from "./common/switch.js";

export const BinaryLight = ({ id }) => {
  const { toggle, on } = getSwitch(id);
  return (
    <IconToggle
      iconOn="lightbulbOn"
      iconOff="lightbulbOnOutline"
      toggle={toggle}
      on={on}
    />
  );
};

BinaryLight.propTypes = {
  id: PropTypes.string.isRequired,
};
