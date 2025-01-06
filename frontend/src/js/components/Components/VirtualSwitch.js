import React from "react";
import PropTypes from "prop-types";

import { mdiLightSwitch, mdiLightSwitchOff } from "@mdi/js";

import { IconToggle } from "./common/IconToggle";
import { getSwitch } from "./common/switch.js";

export const VirtualSwitch = ({ id }) => {
  const { toggle, on } = getSwitch(id);

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
