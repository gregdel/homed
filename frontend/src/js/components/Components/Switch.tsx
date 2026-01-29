import React from "react";
import { useComponents } from "../ComponentsContext";
import { SwitchValues } from "../../types";

import { IconToggle } from "./common/IconToggle";

interface SwitchProps {
  id: string;
}

export const Switch: React.FC<SwitchProps> = ({ id }) => {
  const { getComponentById, updateComponent } = useComponents();

  const component = getComponentById(id);
  if (component === undefined) {
    return null;
  }

  const { on } = component.values as SwitchValues;

  const toggle = () => {
    void updateComponent(id, on ? "OFF" : "ON");
  };

  return <IconToggle iconOn="power" iconOff="power" toggle={toggle} on={on} />;
};
