import React from "react";
import PropTypes from "prop-types";

import { IconToggle } from "./common/IconToggle";
import { getSwitch } from "./common/switch.js";

export const BinaryFan = ({ id }) => {
  const { toggle, on } = getSwitch(id);
  return (
    <IconToggle
      iconOn="fan"
      iconOff="fan"
      toggle={toggle}
      on={on}
      rotate={on}
    />
  );
};

BinaryFan.propTypes = {
  id: PropTypes.string.isRequired,
};
