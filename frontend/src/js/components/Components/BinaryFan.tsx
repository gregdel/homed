import React from "react";

import { IconToggle } from "./common/IconToggle";
import { getSwitch } from "./common/switch";

interface BinaryFanProps {
  id: string;
}

export const BinaryFan: React.FC<BinaryFanProps> = ({ id }) => {
  const result = getSwitch(id);
  if (!result) return null;

  const { toggle, on } = result;
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
