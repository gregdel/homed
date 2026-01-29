import type React from "react";
import type {
  BinarySensorValues,
  GenericSensorValues,
  WifiSignalValues,
} from "../../types";
import { useComponents } from "../ComponentsContext";

interface SensorComponentProps {
  id: string;
}

// Helper to format unknown sensor values
const formatValue = (val: unknown): string => {
  if (val === null || val === undefined) return "—";
  if (
    typeof val === "string" ||
    typeof val === "number" ||
    typeof val === "boolean"
  ) {
    return String(val);
  }
  return JSON.stringify(val);
};

// Generic Sensor Component
export const GenericSensor: React.FC<SensorComponentProps> = ({ id }) => {
  const { getComponentById } = useComponents();

  const component = getComponentById(id);
  if (component === undefined) return null;

  const values = component.values as GenericSensorValues;
  const displayValue = formatValue(values.value);

  return (
    <div style={{ display: "flex", justifyContent: "space-evenly" }}>
      <div style={{ fontSize: "2em", fontWeight: 200 }}>
        <span>{displayValue}</span>
      </div>
    </div>
  );
};

// Power Meter Component (formatted sensor display)
export const PowerMeter: React.FC<SensorComponentProps> = ({ id }) => {
  const { getComponentById } = useComponents();
  const component = getComponentById(id);
  if (component === undefined) return null;

  const values = component.values as GenericSensorValues;

  return (
    <div style={{ fontSize: "4em", fontWeight: 200 }}>
      <span>{String(values.value)} W</span>
    </div>
  );
};

// WiFi Signal Component
export const WifiSignal: React.FC<SensorComponentProps> = ({ id }) => {
  const { getComponentById } = useComponents();
  const component = getComponentById(id);
  if (component === undefined) return null;

  const { value } = component.values as WifiSignalValues;
  if (!value) return null;

  return <>{value}dB</>;
};

// Device Status Component
export const DeviceStatus: React.FC<SensorComponentProps> = ({ id }) => {
  const { getComponentById } = useComponents();

  const component = getComponentById(id);
  if (component === undefined) return null;

  const { on: online } = component.values as BinarySensorValues;
  return <>{online ? "Online" : "Offline"}</>;
};
