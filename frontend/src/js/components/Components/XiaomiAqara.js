import React from "react";
import PropTypes from "prop-types";
import { useSelector } from "react-redux";

export const XiaomiAqara = ({ id }) => {
  const { humidity, pressure, temperature } = useSelector(
    (state) => state.components.components.get(id).values
  );

  return (
    <div>
      <div>Temperature: {temperature}°C</div>
      <div>Humidity: {humidity}%</div>
      <div>Pressure: {pressure}hPa</div>
    </div>
  );
};

XiaomiAqara.propTypes = {
  id: PropTypes.string.isRequired,
};
