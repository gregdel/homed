import React from "react";
import { useComponents } from "../ComponentsContext";
import type { BinarySensorValues } from "../../types";

interface DeviceStatusProps {
  id: string;
}

export const DeviceStatus: React.FC<DeviceStatusProps> = ({ id }) => {
  const { getComponentById } = useComponents();

  const component = getComponentById(id);
  if (component === undefined) {
    return null;
  }

  const { on: online } = component.values as BinarySensorValues;
  return <>{online ? "Online" : "Offline"}</>;
};
