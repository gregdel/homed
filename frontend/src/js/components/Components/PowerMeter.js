import React from "react";
import PropTypes from "prop-types";
import { useComponents } from "../ComponentsContext";

export const PowerMeter = ({ id }) => {
  const { getComponentById } = useComponents();
  const component = getComponentById(id);
  if (component === undefined) {
    return null;
  }

  return (
    <div style={{ fontSize: "4em", fontWeight: 200 }}>
      <span>{component.values.value} W</span>
    </div>
  );
};

PowerMeter.propTypes = {
  id: PropTypes.string.isRequired,
};
