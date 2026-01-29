import React from "react";
import PropTypes from "prop-types";

import { IconToggle } from "./common/IconToggle";
import { getSwitch } from "./common/switch.js";

export const VirtualSwitch = ({ id }) => {
  const { toggle, on } = getSwitch(id);

  return (
    <IconToggle
      iconOn="lightSwitch"
      iconOff="lightSwitchOff"
      toggle={toggle}
      on={on}
    />
  );
};

VirtualSwitch.propTypes = {
  id: PropTypes.string.isRequired,
};
