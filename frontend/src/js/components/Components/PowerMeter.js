import React from "react";
import { useSelector } from "react-redux";
import PropTypes from "prop-types";

export const PowerMeter = ({ id }) => {
  const data = useSelector((state) => state.components.components.get(id));
  if (!data) {
    return null;
  }

  return (
    <div style={{ fontSize: "4em", fontWeight: 200 }}>
      <span>{data.values.value} W</span>
    </div>
  );
};

PowerMeter.propTypes = {
  id: PropTypes.string.isRequired,
};
