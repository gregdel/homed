import React from "react";
import PropTypes from "prop-types";
import { useComponents } from "../ComponentsContext";

export const GenericSensor = ({ id }) => {
  const { getComponentById } = useComponents();

  const component = getComponentById(id);
  if (component === undefined) {
    return null;
  }

  return (
    <div style={{ display: "flex", justifyContent: "space-evenly" }}>
      <div style={{ fontSize: "2em", fontWeight: 200 }}>
        <span>{component.values.value}</span>
      </div>
    </div>
  );
};

GenericSensor.propTypes = {
  id: PropTypes.string.isRequired,
};
