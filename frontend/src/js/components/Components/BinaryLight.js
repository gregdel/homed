import React from "react";
import PropTypes from "prop-types";

import { mdiLightbulbOnOutline, mdiLightbulbOn } from "@mdi/js";

import { IconToggle } from "./common/IconToggle";
import { getSwitch } from "./common/switch.js";

export const BinaryLight = ({ id }) => {
  const { toggle, on } = getSwitch(id);
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
