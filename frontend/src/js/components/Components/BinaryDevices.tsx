import React from "react";
import type { IconName } from "../ui";
import { Icon } from "../ui";
import { IconToggle } from "./common/IconToggle";
import { getSwitch } from "./common/switch";

interface BinaryComponentProps {
  id: string;
}

// Binary Light Component
export const BinaryLight: React.FC<BinaryComponentProps> = ({ id }) => {
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

// Binary Fan Component
export const BinaryFan: React.FC<BinaryComponentProps> = ({ id }) => {
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

// Binary Sensor Component
export const BinarySensor: React.FC<BinaryComponentProps> = ({ id }) => {
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

// Binary TRV Component
export const BinaryTRV: React.FC<BinaryComponentProps> = ({ id }) => {
  const switchData = getSwitch(id);
  if (!switchData) return null;

  const { on } = switchData;

  return (
    <>
      <div style={{ display: "flex", justifyContent: "center" }}>
        <Icon
          name="radiator"
          size={3}
          style={{
            color: on ? "#ff00005e" : "#00000040",
            transition: "color 0.3s ease-out 0s",
          }}
        />
      </div>
      <div style={{ display: "flex", justifyContent: "center" }}>
        <span>TRV is {on ? "open" : "closed"}</span>
      </div>
    </>
  );
};

// Boiler Component
export const Boiler: React.FC<BinaryComponentProps> = ({ id }) => {
  const switchData = getSwitch(id);
  if (!switchData) return null;

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

// Generic binary device component - for internal use
interface GenericBinaryDeviceProps {
  id: string;
  iconOn: IconName;
  iconOff: IconName;
  rotate?: boolean;
}

export const GenericBinaryDevice: React.FC<GenericBinaryDeviceProps> = ({
  id,
  iconOn,
  iconOff,
  rotate = false,
}) => {
  const result = getSwitch(id);
  if (!result) return null;

  const { toggle, on } = result;
  return (
    <IconToggle
      iconOn={iconOn}
      iconOff={iconOff}
      toggle={toggle}
      on={on}
      rotate={rotate}
    />
  );
};
