import React from "react";
import PropTypes from "prop-types";

export const Humidity = ({ value }) => {
  return <>{value}%H</>;
};

Humidity.propTypes = {
  value: PropTypes.number.isRequired,
};
