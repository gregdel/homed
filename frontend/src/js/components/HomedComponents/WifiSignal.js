import React from "react";
import PropTypes from "prop-types";

export const WifiSignal = ({ value }) => {
  return <>{value}dB</>;
};

WifiSignal.propTypes = {
  value: PropTypes.number.isRequired,
};
