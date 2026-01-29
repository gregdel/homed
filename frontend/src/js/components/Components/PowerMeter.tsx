import React from "react";
import { useComponents } from "../ComponentsContext";
import type { GenericSensorValues } from "../../types";

interface PowerMeterProps {
  id: string;
}

export const PowerMeter: React.FC<PowerMeterProps> = ({ id }) => {
  const { getComponentById } = useComponents();
  const component = getComponentById(id);
  if (component === undefined) {
    return null;
  }

  const values = component.values as GenericSensorValues;

  return (
    <div style={{ fontSize: "4em", fontWeight: 200 }}>
      <span>{String(values.value)} W</span>
    </div>
  );
};
