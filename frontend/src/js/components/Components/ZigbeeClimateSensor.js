import React from "react";
import PropTypes from "prop-types";
import { useComponents } from "../ComponentsContext";

import { Icon } from "../ui/Icon";

export const ZigbeeClimateSensor = ({ id }) => {
  const { getComponentById } = useComponents();
  const component = getComponentById(id);
  if (component === undefined) {
    return null;
  }

  const { humidity, pressure, temperature } = component.values;

  const prettyTemperature = (t) => t.toFixed(2);

  return (
    <div>
      <div style={{ display: "flex", marginBottom: "0.3em" }}>
        <Icon name="thermometer" size={1} />
        Temperature: {prettyTemperature(temperature)}°C
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

ZigbeeClimateSensor.propTypes = {
  id: PropTypes.string.isRequired,
};
