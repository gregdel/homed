import React from "react";

import { Icon } from "../ui/Icon";
import { getSwitch } from "./common/switch";

interface BoilerProps {
  id: string;
}

export const Boiler: React.FC<BoilerProps> = ({ id }) => {
  const switchData = getSwitch(id);

  if (!switchData) {
    return null;
  }

  const { on } = switchData;

  return (
    <div
      style={{
        display: "flex",
        alignItems: "flex-end",
        flexWrap: "wrap",
        justifyContent: "space-between",
      }}
    >
      <div>
        <h2>Boiler</h2>
      </div>
      <div>
        <div>
          <Icon
            name="fire"
            size={3}
            style={{
              color: on ? "#ff4d4f" : "#00000040",
              transition: "color 0.3s ease-out 0s",
            }}
          />
        </div>
      </div>
    </div>
  );
};
