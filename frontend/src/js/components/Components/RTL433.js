import React from "react";
import PropTypes from "prop-types";
import { useSelector } from "react-redux";

export const RTL433 = ({ id }) => {
  const state = useSelector(
    (state) => state.components.components.get(id).values.on
  );
  return <>{state ? "Boiler is on" : "Boiler is off"}</>;
};

RTL433.propTypes = {
  id: PropTypes.string.isRequired,
};
