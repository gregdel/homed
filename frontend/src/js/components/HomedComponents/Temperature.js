import React from "react";
import PropTypes from "prop-types";

export const Temperature = ({ value }) => {
  return <>{value}°C</>;
};

Temperature.propTypes = {
  value: PropTypes.number.isRequired,
};
