import React from "react";

import { Icon } from "../ui/Icon";
import { getSwitch } from "./common/switch";

interface BinarySensorProps {
  id: string;
}

export const BinarySensor: React.FC<BinarySensorProps> = ({ id }) => {
  const result = getSwitch(id);
  if (!result) return null;

  const { on } = result;

  return (
    <div style={{ display: "flex", justifyContent: "space-evenly" }}>
      <Icon name={on ? "toggleSwitch" : "toggleSwitchOffOutline"} size={2} />
      <div style={{ fontSize: "2em", fontWeight: 200 }}>
        <span>{on ? "ON" : "OFF"}</span>
      </div>
    </div>
  );
};
