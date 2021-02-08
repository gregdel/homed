import React from "react";
import PropTypes from "prop-types";
import { useSelector } from "react-redux";

export const Humidity = ({ id }) => {
  const value = useSelector(
    (state) => state.components.components.get(id).values.value
  );
  return <>{value}%H</>;
};

Humidity.propTypes = {
  id: PropTypes.string.isRequired,
};
