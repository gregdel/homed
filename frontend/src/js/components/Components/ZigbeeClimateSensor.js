import React from "react";
import PropTypes from "prop-types";
import { useSelector } from "react-redux";

import Icon from "@mdi/react";
import { mdiThermometer, mdiWaterPercent, mdiGauge } from "@mdi/js";

export const ZigbeeClimateSensor = ({ id }) => {
  const { humidity, pressure, temperature } = useSelector(
    (state) => state.components.components.get(id).values
  );

  return (
    <div>
      <div style={{ display: "flex", marginBottom: "0.3em" }}>
        <Icon path={mdiThermometer} size={1} />
        Temperature: {temperature}°C
      </div>
      <div style={{ display: "flex", marginBottom: "0.3em" }}>
        <Icon path={mdiWaterPercent} size={1} />
        Humidity: {humidity}%
      </div>
      {pressure !== 0 && (
        <div style={{ display: "flex" }}>
          <Icon path={mdiGauge} size={1} />
          Pressure: {pressure}hPa
        </div>
      )}
    </div>
  );
};

ZigbeeClimateSensor.propTypes = {
  id: PropTypes.string.isRequired,
};
