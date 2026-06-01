import type React from "react";
import { useComponents } from "../ComponentsContext";

import { Icon } from "../ui/Icon";

interface ClimateSensorProps {
  id: string;
}

interface ClimateSensorValues {
  humidity: number;
  pressure: number;
  temperature: number;
  [key: string]: unknown;
}

export const ClimateSensor: React.FC<ClimateSensorProps> = ({ id }) => {
  const { getComponentById } = useComponents();
  const component = getComponentById(id);
  if (component === undefined) {
    return null;
  }

  const { humidity, pressure, temperature } =
    component.values as unknown as ClimateSensorValues;

  const prettyTemperature = (t: number) => t.toFixed(2);

  return (
    <div>
      <div style={{ display: "flex", marginBottom: "0.3em" }}>
        <Icon name="thermometer" size={1} />
        Temperature: {prettyTemperature(temperature)}
        °C
      </div>
      <div style={{ display: "flex", marginBottom: "0.3em" }}>
        <Icon name="waterPercent" size={1} />
        Humidity: {humidity}%
      </div>
      {pressure !== 0 && (
        <div style={{ display: "flex" }}>
          <Icon name="gauge" size={1} />
          Pressure: {pressure}hPa
        </div>
      )}
    </div>
  );
};
