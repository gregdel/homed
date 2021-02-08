import React from "react";
import PropTypes from "prop-types";
import { useSelector } from "react-redux";

export const Temperature = ({ id }) => {
  const value = useSelector(
    (state) => state.components.components.get(id).values.value
  );
  return <>{value}°C</>;
};

Temperature.propTypes = {
  id: PropTypes.string.isRequired,
};
