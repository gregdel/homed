import React from "react";
import { useComponents } from "../ComponentsContext";
import type { WifiSignalValues } from "../../types";

interface WifiSignalProps {
  id: string;
}

export const WifiSignal: React.FC<WifiSignalProps> = ({ id }) => {
  const { getComponentById } = useComponents();
  const component = getComponentById(id);
  if (component === undefined) {
    return null;
  }

  const { value } = component.values as WifiSignalValues;
  if (!value) {
    return null;
  }

  return <>{value}dB</>;
};
