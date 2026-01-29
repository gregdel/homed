import React from "react";
import { Icon } from "../ui/Icon";
import { getSwitch } from "./common/switch";

interface BinaryTRVProps {
  id: string;
}

export const BinaryTRV: React.FC<BinaryTRVProps> = ({ id }) => {
  const switchData = getSwitch(id);
  if (!switchData) {
    return null;
  }

  const { on } = switchData;

  return (
    <>
      <div
        style={{
          display: "flex",
          justifyContent: "center",
        }}
      >
        <Icon
          name="radiator"
          size={3}
          style={{
            color: on ? "#ff00005e" : "#00000040",
            transition: "color 0.3s ease-out 0s",
          }}
        />
      </div>
      <div
        style={{
          display: "flex",
          justifyContent: "center",
        }}
      >
        <span>TRV is {on ? "open" : "closed"}</span>
      </div>
    </>
  );
};
