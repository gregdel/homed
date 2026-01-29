import React from "react";
import { useComponents } from "../ComponentsContext";
import type { SwitchValues } from "../../types";
import { IconToggle } from "./common/IconToggle";

interface SwitchComponentProps {
  id: string;
}

// Regular Switch Component
export const Switch: React.FC<SwitchComponentProps> = ({ id }) => {
  const { getComponentById, updateComponent } = useComponents();

  const component = getComponentById(id);
  if (component === undefined) return null;

  const { on } = component.values as SwitchValues;

  const toggle = () => {
    void updateComponent(id, on ? "OFF" : "ON");
  };

  return <IconToggle iconOn="power" iconOff="power" toggle={toggle} on={on} />;
};

// Virtual Switch Component (uses getSwitch helper)
export const VirtualSwitch: React.FC<SwitchComponentProps> = ({ id }) => {
  const { getComponentById, updateComponent } = useComponents();

  const component = getComponentById(id);
  if (component === undefined) return null;

  const { on } = component.values as SwitchValues;

  const toggle = () => {
    void updateComponent(id, on ? "OFF" : "ON");
  };

  return (
    <IconToggle
      iconOn="lightSwitch"
      iconOff="lightSwitchOff"
      toggle={toggle}
      on={on}
    />
  );
};
