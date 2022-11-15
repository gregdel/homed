import React from "react";
import { useSelector } from "react-redux";
import PropTypes from "prop-types";

export const GenericSensor = ({ id }) => {
  const data = useSelector((state) => state.components.components.get(id));
  if (!data) {
    return null;
  }

  return (
    <div style={{ display: "flex", justifyContent: "space-evenly" }}>
      <div style={{ fontSize: "2em", fontWeight: 200 }}>
        <span>{data.values.value}</span>
      </div>
    </div>
  );
};

GenericSensor.propTypes = {
  id: PropTypes.string.isRequired,
};
