import React from "react";
import PropTypes from "prop-types";
import { useSelector } from "react-redux";

export const WifiSignal = ({ id }) => {
  const value = useSelector(
    (state) => state.components.components.get(id).values.value
  );
  if (!value) {
    return null;
  }

  return <>{value}dB</>;
};

WifiSignal.propTypes = {
  id: PropTypes.string.isRequired,
};
