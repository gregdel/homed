import React from "react";
import { useComponents } from "../ComponentsContext";

import { Icon } from "../ui/Icon";

interface ZigbeeClimateSensorProps {
  id: string;
}

interface ZigbeeClimateSensorValues {
  humidity: number;
  pressure: number;
  temperature: number;
  [key: string]: unknown;
}

export const ZigbeeClimateSensor: React.FC<ZigbeeClimateSensorProps> = ({
  id,
}) => {
  const { getComponentById } = useComponents();
  const component = getComponentById(id);
  if (component === undefined) {
    return null;
  }

  const { humidity, pressure, temperature } =
    component.values as unknown as ZigbeeClimateSensorValues;

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
