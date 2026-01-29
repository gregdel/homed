import React from "react";

import { IconToggle } from "./common/IconToggle";
import { getSwitch } from "./common/switch";

interface VirtualSwitchProps {
  id: string;
}

export const VirtualSwitch: React.FC<VirtualSwitchProps> = ({ id }) => {
  const result = getSwitch(id);
  if (!result) return null;

  const { toggle, on } = result;

  return (
    <IconToggle
      iconOn="lightSwitch"
      iconOff="lightSwitchOff"
      toggle={toggle}
      on={on}
    />
  );
};
