import React from "react";

import { IconToggle } from "./common/IconToggle";
import { getSwitch } from "./common/switch";

interface BinaryLightProps {
  id: string;
}

export const BinaryLight: React.FC<BinaryLightProps> = ({ id }) => {
  const result = getSwitch(id);
  if (!result) return null;

  const { toggle, on } = result;
  return (
    <IconToggle
      iconOn="lightbulbOn"
      iconOff="lightbulbOnOutline"
      toggle={toggle}
      on={on}
    />
  );
};
