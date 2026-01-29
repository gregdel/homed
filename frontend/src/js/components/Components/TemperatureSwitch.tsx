import React from "react";
import { useComponents } from "../ComponentsContext";
import type { SwitchValues } from "../../types";

import { Icon } from "../ui/Icon";

interface TemperatureSwitchProps {
  id: string;
}

export const TemperatureSwitch: React.FC<TemperatureSwitchProps> = ({ id }) => {
  const { getComponentById, updateComponent } = useComponents();

  const component = getComponentById(id);
  if (component === undefined) {
    return null;
  }

  const { on } = component.values as SwitchValues;

  const toggle = () => {
    void updateComponent(id, on ? "OFF" : "ON");
  };

  return (
    <div onClick={toggle}>
      <Icon
        name={on ? "thermometer" : "thermometerOff"}
        size={3}
        style={{
          cursor: "pointer",
          color: on ? "#000000" : "#00000040",
          transition: "color 0.3s ease-out 0s",
        }}
      />
    </div>
  );
};
