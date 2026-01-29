import React from "react";
import { useComponents } from "../ComponentsContext";
import type { GenericSensorValues } from "../../types";

interface GenericSensorProps {
  id: string;
}

export const GenericSensor: React.FC<GenericSensorProps> = ({ id }) => {
  const { getComponentById } = useComponents();

  const component = getComponentById(id);
  if (component === undefined) {
    return null;
  }

  const values = component.values as GenericSensorValues;

  const formatValue = (val: unknown): string => {
    if (val === null || val === undefined) {
      return "—";
    }
    if (
      typeof val === "string" ||
      typeof val === "number" ||
      typeof val === "boolean"
    ) {
      return String(val);
    }
    return JSON.stringify(val);
  };

  const displayValue = formatValue(values.value);

  return (
    <div style={{ display: "flex", justifyContent: "space-evenly" }}>
      <div style={{ fontSize: "2em", fontWeight: 200 }}>
        <span>{displayValue}</span>
      </div>
    </div>
  );
};
